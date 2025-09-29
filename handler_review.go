package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
)

const ghRegex = `^(?:https://github\.com/|git@github\.com:)([^/]+)/([^/]+?)(?:\.git)?/?$`

func isValidGitHubURL(link string) bool {
	var regex = regexp.MustCompile(ghRegex)
	return regex.MatchString(link)

}

func extractRepoDetails(link string) (string, string) {
	var regex = regexp.MustCompile(ghRegex)
	matches := regex.FindStringSubmatch(link)

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

	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			errorOutput := stderr.String()
			if strings.Contains(errorOutput, "not found") || strings.Contains(errorOutput, "could not read") {
				return fmt.Errorf("Repository not found or is private. Please check:\n  - URL is correct\n  - Repository is set to public\n")
			}
			return fmt.Errorf("git clone failed (exit %d)", ee.ExitCode())
		}
		return fmt.Errorf("failed to run git: %w", err)
	}

	return nil
}

func deleteDirectory(path string) error {
	cmd := exec.Command("rm", "-rf", path)
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("failed to delete local repo (exit %d)", ee.ExitCode())
		}
		return fmt.Errorf("failed to delete directory: %w", err)
	}
	fmt.Println("✅ Directory deleted...")
	return nil

}

func setupLocalRepo(repositoryLink, localRepoPath string) (string, error) {
	for dirExists(localRepoPath) {
		fmt.Println("\nDirectory already exists. Choose an option:")
		fmt.Println("  [Enter] - Delete and overwrite")
		fmt.Println("  [n]     - Clone with different name")
		fmt.Println("  [q]     - Cancel operation")
		fmt.Print("Your choice: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		if strings.ToLower(input.Text()) == "" {
			deleteDirectory(localRepoPath)
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
		return "", fmt.Errorf("❌ Could not clone repo. \n%v", err)
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

	fmt.Println("\nStarting dev server...")
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
		fmt.Println("\nDev server stopped.")

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
			fmt.Printf("\nCould not start file server: %v", err)
		}
	}()

	openHTMLFile("http://localhost:5500")
	return nil
}

func openHTMLFile(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	default:
		cmd = exec.Command("xdg-open", filePath)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Could not open browser: %v", err)
	}
	return nil
}

func cleanUp(dirPath string) error {
	for {
		fmt.Println("\nPress Enter to delete project directory, or type 's' to save it:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if input == "" {
			if err := deleteDirectory(dirPath); err != nil {
				return fmt.Errorf("Cleanup failed: %v", err)
			}
			fmt.Println("🗑️ Cleanup successful...")
			return nil
		} else if input == "s" {
			fmt.Println("Directory saved. Don't forget to remove it when done")
			return nil
		} else {
			fmt.Println("Invalid input. Press Enter or type 's'...")
		}
	}
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
		return err
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

	err = cleanUp(localRepoPath)
	if err != nil {
		return fmt.Errorf("Could not delete directory")
	}
	fmt.Println("Shutting down.. Goodbye!")
	return nil
}
