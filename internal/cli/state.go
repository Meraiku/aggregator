package cli

import (
	"github.com/meraiku/aggregator/internal/config"
	"github.com/meraiku/aggregator/internal/repository/sql"
)

type State struct {
	cfg *config.Config
	db  *sql.Queries
}

func NewState(cfg *config.Config) (*State, error) {

	db, err := sql.ConnectPostgres(cfg.DbURL)
	if err != nil {
		return nil, err
	}

	return &State{
		cfg: cfg,
		db:  db,
	}, nil
}
