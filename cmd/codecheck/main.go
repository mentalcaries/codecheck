package main

import (
	"log"
	"os"

	"github.com/mentalcaries/codecheck/internal/config"
)

const Version = "v0.1.0"

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
	cmds.register("setup", handlerSetup)

	commandArgs := os.Args
	if len(commandArgs) < 2 {
		log.Fatal("Invalid command. Usage: codecheck <command>")
	}

	cmd := command{name: commandArgs[1], args: commandArgs[2:]}

	err = cmds.run(appState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
