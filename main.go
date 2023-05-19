package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	RepoPath     string
	DownloadPath string
}

func LoadConfig(path string) Config {
	config := Config{
		RepoPath:     os.ExpandEnv("$HOME/dev/git/overleaf-archive"),
		DownloadPath: os.ExpandEnv("$HOME/Downloads"),
	}
	// TODO
	return config
}

func (config Config) Save(path string) {
	// TODO
}

func (config Config) Projects() []string {
	dirs, err := os.ReadDir(config.RepoPath)
	if err != nil {
		panic(err)
	}

	projects := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		projects = append(projects, dir.Name())
	}

	return projects
}

func (config Config) DownloadedProjects() []string {
	projects := config.Projects()

	dirs, err := os.ReadDir(config.DownloadPath)
	if err != nil {
		panic(err)
	}

	update_projects := make([]string, 0, len(projects))
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		p := dir.Name()
		for _, project := range projects {
			if p == project {
				update_projects = append(update_projects, dir.Name())
			}
		}
	}
	return update_projects
}

func main() {
	config_file_path := os.ExpandEnv("$HOME/.overleaf-backup.config.json")

	config := LoadConfig(config_file_path)

	projects := config.DownloadedProjects()
	if len(projects) == 0 {
		return
	}

	fmt.Printf("Update [%d] projects\n", len(projects))
	for i, project := range projects {
		fmt.Printf("[%02d/%02d] %s\n", i+1, len(projects), project)
		fmt.Printf("commit> ")
		message := ""
		fmt.Scanf("%s", &message)
		message = strings.TrimSpace(message)
		if len(message) == 0 {
			message = project
		} else {
			message = message + ": " + project
		}

		exec.Command("rm", "-r", config.RepoPath+"/"+project).Run()
		exec.Command("mv", config.DownloadPath+"/"+project, config.RepoPath+"/"+project).Run()
		os.Chdir(config.RepoPath)
		exec.Command("git", "add", project).Run()
		exec.Command("git", "commit", "-m", message).Run()
	}

	os.Chdir(config.RepoPath)
	exec.Command("git", "push").Run()
}
