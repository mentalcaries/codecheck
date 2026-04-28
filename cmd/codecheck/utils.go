package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

const ghRegex = `^(?:https://github\.com/|git@github\.com:)([^/]+)/([^/]+?)(?:\.git)?(?:/tree/([^/]+))?/?$`

const ttProjectRegex = `^([a-zA-Z0-9-]+)-(se_project_[a-zA-Z0-9_]+)-([0-9a-f]{7,40})$`

func isValidGitHubURL(link string) bool {
	var regex = regexp.MustCompile(ghRegex)
	return regex.MatchString(link)
}

func isValidProjectString(link string) bool {
	var regex = regexp.MustCompile(ttProjectRegex)
	return regex.MatchString(link)
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func fileExists(dirPath, filename string) bool {
	pathToFile := filepath.Join(dirPath, filename)
	_, err := os.Stat(pathToFile)
	return err == nil
}

func hasPackageJSON(dirPath string) bool {
	packagePath := filepath.Join(dirPath, "package.json")
	_, err := os.Stat(packagePath)
	return err == nil
}

func hasIndexHTML(dirPath string) bool {
	packagePath := filepath.Join(dirPath, "index.html")
	_, err := os.Stat(packagePath)
	return err == nil
}

func getProjectRoot(dirPath string) (string, error) {
	var rootPath string
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && rootPath == "" {
			if hasPackageJSON(dirPath) || hasIndexHTML(dirPath) {
				rootPath = path
				return filepath.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return rootPath, nil
}
