package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
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
	return err == nil
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

func setupLocalRepo(repositoryLink, localRepoPath string) (string, error) {
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
			return "", fmt.Errorf("User cancelled operation")
		}
	}

	err := cloneRepository(repositoryLink, localRepoPath)
	if err != nil {
		return "", fmt.Errorf("Could not clone repo: %v", err)
	}
	return localRepoPath, nil
}

func isVSCodeInstalled() bool {
	_, err := exec.LookPath("code")
	return err == nil
}

func openRepoWithVSCode(repoPath string) error {
	if isVSCodeInstalled() {
		fmt.Println("VS Code available. Opening project...")
		cmd := exec.Command("code", repoPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Error opening VSCode")
		}
	} else {
		fmt.Println("VS Code CLI not found")
	}

	return nil
}

func hasPackageJSON(dirPath string) bool {
	packagePath := filepath.Join(dirPath, "package.json")
	_, err := os.Stat(packagePath)
	return err == nil
}

func installDependencies(dirPath string) error {
	cmd := exec.Command("npm", "install")
	cmd.Dir = dirPath
	cmd.Stderr = os.Stderr

	fmt.Println("Installing NPM depedencies...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Could not install dependencies: %w", err)
	} else {
		fmt.Println("Depedencies installed successfully.")
	}

	return nil
}

func startDevServer(dirPath string) error {
	cmd := exec.Command("npm", "run", "dev")
	cmd.Dir = dirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("\n>>>Starting dev server...")
	if err := cmd.Start(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("Failed to start the dev server (exit %d)", ee.ExitCode())
		}
		return fmt.Errorf("Could not start dev server: %w", err)
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("\n>>>Dev server stopped.")

		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()
	return nil
}

func startFileServer(dirPath string) error {
	fileserver := http.FileServer(http.Dir(dirPath))
	http.Handle("/", fileserver)

	fmt.Println("Serving index.html on Port 5500. Access the project here: http://localhost:5500")

	go func() {
		if err := http.ListenAndServe(":5500", nil); err != nil {
			fmt.Errorf("Could not start file server: %v", err)
		}
	}()
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

	localRepoPath, err = setupLocalRepo(repositoryLink, localRepoPath)
	if err != nil {
		return fmt.Errorf("Could not set up local directory: %v", err)
	}

	if hasPackageJSON(localRepoPath) {
		fmt.Println("\nReading package.json...")
		if err = installDependencies(localRepoPath); err != nil {
			return fmt.Errorf("%v", err)
		}
		if err = startDevServer(localRepoPath); err != nil {
			return fmt.Errorf("%v", err)
		}
	} else {
		if err = startFileServer(localRepoPath); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	openRepoWithVSCode(localRepoPath)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Shutting down.. Goodbye!")
	return nil
}
