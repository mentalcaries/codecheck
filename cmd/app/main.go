package main

import (
	"fmt"
	"log"

	"github.com/mentalcaries/tt-codecheck/internal/config"
)

type state struct {
  config *config.Config
}

func main() {
  cfg, err := config.Read()

  if err != nil {
    log.Fatalf("error reading config: %v", err)
  }

  appState := &state{
    config: &cfg,
  }

  fmt.Println(appState.config.DownloadDirectory)
}