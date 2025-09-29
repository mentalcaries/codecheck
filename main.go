package main

import (
	"log"
	"os"

	"github.com/mentalcaries/tt-codecheck/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.CheckConfig()

	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	appState := &state{
		config: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("review", handlerReview)

	commandArgs := os.Args
	if len(commandArgs) < 2 {
		log.Fatal("Invalid command. Usage: codecheck <github-repo-url>")
	}

	cmd := command{name: commandArgs[1], args: commandArgs[2:]}

	err = cmds.run(appState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
