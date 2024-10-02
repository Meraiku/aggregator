package main

import (
	"fmt"
	"log"

	"github.com/meraiku/aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	fmt.Println(cfg)

	err = cfg.SetUser("Meraiku")
	if err != nil {
		log.Fatalf("setting user in config: %s", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	fmt.Println(cfg)

}
