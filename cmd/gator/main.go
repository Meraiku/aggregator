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

	err = cmds.RegisterHandlers()
	if err != nil {
		log.Fatalf("error register handlers: %s", err)
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
