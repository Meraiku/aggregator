package main

import (
	"log"
	"os"

	"github.com/meraiku/aggregator/internal/cli"
	"github.com/meraiku/aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	state, err := cli.NewState(cfg)
	if err != nil {
		log.Fatalf("error creating new state: %s", err)
	}

	cmds := cli.NewCommands()

	err = cmds.Register("login", cli.Login)
	if err != nil {
		log.Fatalf("error register login handler: %s", err)
	}

	err = cmds.Register("register", cli.Register)
	if err != nil {
		log.Fatalf("error register register handler: %s", err)
	}

	err = cmds.Register("reset", cli.Reset)
	if err != nil {
		log.Fatalf("error register reset handler: %s", err)
	}

	err = cmds.Register("users", cli.Users)
	if err != nil {
		log.Fatalf("error register users handler: %s", err)
	}

	args := os.Args
	if len(args) < 2 {
		log.Fatal(cli.ErrInvalidArgumentsCount)
	}

	cmd := cli.NewCommand(args[1], args[2:])

	if err := cmds.Run(state, *cmd); err != nil {
		log.Fatal(err)
	}
}
