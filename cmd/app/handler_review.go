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

	tempDirPath := filepath.Join(s.config.DownloadDirectory, userName)

	fmt.Println("temp dir: ", tempDirPath)

	if !dirExists(tempDirPath) {
		fmt.Println("dir does not exist")
	}


	return nil
}
