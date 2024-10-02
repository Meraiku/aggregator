package cli

import "github.com/meraiku/aggregator/internal/config"

type State struct {
	cfg *config.Config
}

func NewState(cfg *config.Config) *State {
	return &State{
		cfg: cfg,
	}
}
