package main

import (
	"fmt"

	"github.com/mentalcaries/tt-codecheck/internal/config"
)

func handlerSetup(s *state, cmd command) error {
	fmt.Print("Enter download directory [current: ")

	// Show current config if it exists
	if config.ConfigExists() {
		currentConfig, _ := config.Read()
		fmt.Printf("%s]: ~/", currentConfig.DownloadDirectory)
	} else {
		fmt.Print("~/]: ")
	}

	config.ConfigDownloadDir()
	fmt.Println("✅ Config updated successfully!")
    updatedConfig, _ := config.Read()
    fmt.Println("Projects will now be cloned to: ", updatedConfig.DownloadDirectory)
	return nil
}
