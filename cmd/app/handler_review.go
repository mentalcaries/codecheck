package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
		fmt.Println("directory does not exist, creating directory...")
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return "", fmt.Errorf("Could not create user directory")
		}
	}
	fmt.Println("directory exists, using it...")

	return path, nil
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

	userName, _ := extractRepoDetails(repositoryLink)

	userDirPath := filepath.Join(s.config.DownloadDirectory, userName)

	userDirPath, err := setupTempDir(userDirPath)
	if err != nil {
		return fmt.Errorf("Could not create temp directory: %v", err)
	}

	return nil
}
