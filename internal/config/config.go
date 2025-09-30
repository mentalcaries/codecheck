package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	DownloadDirectory string `json:"download_directory"`
}

const configFileName = ".codecheckconfig.json"

func (cfg *Config) SetDownloadDir(dirName string) error {
	cfg.DownloadDirectory = dirName
	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	homeLocation, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}
	configPath := filepath.Join(homeLocation, configFileName)

	return configPath, nil
}

func ConfigExists() bool {
	configPath, err := getConfigFilePath()
	if err != nil {
		fmt.Println("invalid path")
	}
	_, err = os.Stat(configPath)
	return err == nil
}

func ConfigDownloadDir() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	path := scanner.Text()
	if strings.TrimSpace(path) != "" {
		write(Config{DownloadDirectory: path})
	}
}

func CheckConfig() (Config, error) {
	for !ConfigExists() {
		fmt.Printf("Download directory not configured. Enter your temporary download location: ~/")
		ConfigDownloadDir()
	}

	return Read()
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()

	if err != nil {
		fmt.Println("Could not get config path")
		return Config{}, err
	}
	data, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("Could not read config")
		return Config{}, err
	}

	config := Config{}

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("could not get config object", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	config.DownloadDirectory = filepath.Clean(filepath.Join(home, config.DownloadDirectory))

	return config, nil
}

func write(config Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return errors.New("invalid path")
	}

	configJson, _ := json.MarshalIndent(config, "", "  ")

	err = os.WriteFile(configPath, configJson, 0644)
	if err != nil {
		return err
	}

	return nil
}
