package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func isValidGitHubURL(link string) bool {
	var ghRegex = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+?)(?:\.git)?/?$`)
	return ghRegex.MatchString(link)

}

func extractRepoDetails(link string) (string, string) {
	var ghRegex = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+?)(?:\.git)?/?$`)
	matches := ghRegex.FindStringSubmatch(link)

	return matches[1], matches[2]
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func setupTempDir(path string) (string, error) {
	if !dirExists(path) {
		fmt.Println("Directory does not exist, creating directory...")
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return "", fmt.Errorf("Could not create user directory")
		}
	}

	return path, nil
}

func cloneRepository(link, path string) error {
	cmd := exec.Command("git", "clone", link, path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("git clone failed (exit %d)", ee.ExitCode())
		}
		return fmt.Errorf("failed to run git: %w", err)
	}
	return nil
}

func deleteExistingRepo(path string) error {
	cmd := exec.Command("rm", "-rf", path)
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("failed to delete local repo (exit %d)", ee.ExitCode())
		}
		return fmt.Errorf("failed to delete directory: %w", err)
	}
	fmt.Println("Existing directory deleted...")
	return nil

}

func handlerReview(s *state, cmd command) error {

	if len(cmd.args) < 1 {
		return fmt.Errorf("link to github repository is required")
	}

	repositoryLink := cmd.args[0]

	isValidLink := isValidGitHubURL(repositoryLink)
	if !isValidLink {
		return fmt.Errorf("Invalid github URL")
	}

	userName, repoName := extractRepoDetails(repositoryLink)

	userDirPath := filepath.Join(s.config.DownloadDirectory, userName)

	userDirPath, err := setupTempDir(userDirPath)
	if err != nil {
		return fmt.Errorf("Could not create temp directory: %v", err)
	}

	localRepoPath := filepath.Join(userDirPath, repoName)
	for dirExists(localRepoPath) {
		fmt.Print("Repo already exists locally. Overwrite? y/n/q: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		if strings.ToLower(input.Text()) == "y" {
			deleteExistingRepo(localRepoPath)
		}

		if strings.ToLower(input.Text()) == "n" {
			localRepoPath = localRepoPath + "-copy"
			for dirExists(localRepoPath) {
				localRepoPath = localRepoPath + "-copy"
			}
		}

		if strings.ToLower(input.Text()) == "q" {
			return fmt.Errorf("User cancelled operation")
		}

	}

	err = cloneRepository(repositoryLink, localRepoPath)

	return nil
}
