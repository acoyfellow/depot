package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

type Repo struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Branch string `json:"branch"`
	Path   string `json:"path"`
}

type Config struct {
	Repos []Repo `json:"repos"`
}

var (
	configPath   string
	vendorDir    string
)

func init() {
	home := os.Getenv("HOME")
	vendorDir = filepath.Join(home, ".vendor")
	configPath = filepath.Join(vendorDir, "config.json")
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Repos: []Repo{}}, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func saveConfig(config *Config) error {
	// Ensure vendor directory exists
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func cmdAdd(c *cli.Context) error {
	if c.NArg() < 2 {
		return cli.Exit("Usage: depo add <name> <url> [--branch main] [--path ~/.vendor/<name>]", 1)
	}

	name := c.Args().Get(0)
	url := c.Args().Get(1)
	branch := c.String("branch")
	path := c.String("path")

	if branch == "" {
		branch = "main"
	}
	if path == "" {
		path = filepath.Join(vendorDir, name)
	}

	config, err := loadConfig()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error loading config: %v", err), 1)
	}

	// Check if repo already exists
	for _, repo := range config.Repos {
		if repo.Name == name {
			return cli.Exit(fmt.Sprintf("Repo '%s' already exists", name), 1)
		}
	}

	config.Repos = append(config.Repos, Repo{
		Name:   name,
		URL:    url,
		Branch: branch,
		Path:   path,
	})

	if err := saveConfig(config); err != nil {
		return cli.Exit(fmt.Sprintf("Error saving config: %v", err), 1)
	}

	fmt.Printf("✓ Added repo '%s' to config\n", name)
	return nil
}

func cmdRemove(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.Exit("Usage: depo remove <name>", 1)
	}

	name := c.Args().Get(0)
	config, err := loadConfig()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error loading config: %v", err), 1)
	}

	for i, repo := range config.Repos {
		if repo.Name == name {
			config.Repos = append(config.Repos[:i], config.Repos[i+1:]...)
			if err := saveConfig(config); err != nil {
				return cli.Exit(fmt.Sprintf("Error saving config: %v", err), 1)
			}
			fmt.Printf("✓ Removed repo '%s' from config\n", name)
			return nil
		}
	}

	return cli.Exit(fmt.Sprintf("Repo '%s' not found", name), 1)
}

func cmdUpdate(c *cli.Context) error {
	config, err := loadConfig()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error loading config: %v", err), 1)
	}

	if len(config.Repos) == 0 {
		fmt.Println("No repos configured")
		return nil
	}

	name := ""
	if c.NArg() > 0 {
		name = c.Args().Get(0)
	}

	for _, repo := range config.Repos {
		if name != "" && repo.Name != name {
			continue
		}

		fmt.Printf("Processing %s...\n", repo.Name)

		// Check if repo exists
		gitDir := filepath.Join(repo.Path, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			// Update existing repo
			fmt.Println("  Updating existing repo...")
			if err := runGit(repo.Path, "fetch", "origin"); err != nil {
				fmt.Printf("  Error fetching: %v\n", err)
				continue
			}
			if err := runGit(repo.Path, "checkout", repo.Branch); err != nil {
				// Try creating the branch
				if err := runGit(repo.Path, "checkout", "-b", repo.Branch, fmt.Sprintf("origin/%s", repo.Branch)); err != nil {
					fmt.Printf("  Error checking out branch: %v\n", err)
					continue
				}
			}
			if err := runGit(repo.Path, "pull", "origin", repo.Branch); err != nil {
				fmt.Printf("  Error pulling: %v\n", err)
				continue
			}
		} else {
			// Clone new repo
			fmt.Println("  Cloning repo...")
			parent := filepath.Dir(repo.Path)
			if err := os.MkdirAll(parent, 0755); err != nil {
				fmt.Printf("  Error creating directory: %v\n", err)
				continue
			}
			cmd := exec.Command("git", "clone", "--depth", "1", "--branch", repo.Branch, repo.URL, repo.Path)
			if err := cmd.Run(); err != nil {
				fmt.Printf("  Error cloning: %v\n", err)
				continue
			}
		}

		fmt.Printf("  ✓ %s done\n\n", repo.Name)

		if name != "" {
			break
		}
	}

	return nil
}

func cmdList(c *cli.Context) error {
	config, err := loadConfig()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error loading config: %v", err), 1)
	}

	if len(config.Repos) == 0 {
		fmt.Println("No repos configured")
		return nil
	}

	fmt.Println("Configured repos:")
	for _, repo := range config.Repos {
		status := "not cloned"
		if _, err := os.Stat(filepath.Join(repo.Path, ".git")); err == nil {
			status = "cloned"
		}
		fmt.Printf("  %s (%s) - %s [%s]\n", repo.Name, status, repo.URL, repo.Branch)
	}
	return nil
}

func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	app := &cli.App{
		Name:  "depo",
		Usage: "Manage reference repositories globally",
		Commands: []*cli.Command{
			{
				Name:   "add",
				Usage:  "Add a new repo to manage",
				Action: cmdAdd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "branch",
						Value: "main",
						Usage: "Git branch to track",
					},
					&cli.StringFlag{
						Name:  "path",
						Usage: "Local path for repo (default: ~/.vendor/<name>)",
					},
				},
			},
			{
				Name:   "remove",
				Usage:  "Remove a repo from management",
				Action: cmdRemove,
			},
			{
				Name:   "update",
				Usage:  "Update repos (all or specific one)",
				Action: cmdUpdate,
			},
			{
				Name:   "list",
				Usage:  "List configured repos",
				Action: cmdList,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
