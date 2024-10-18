package app

import (
	"github.com/meraiku/aggregator/internal/config"
	"github.com/meraiku/aggregator/internal/repository/sql"
)

type State struct {
	Cfg *config.Config
	Db  *sql.Queries
}

func NewState(cfg *config.Config) (*State, error) {

	db, err := sql.ConnectPostgres(cfg.DbURL)
	if err != nil {
		return nil, err
	}

	return &State{
		Cfg: cfg,
		Db:  db,
	}, nil
}
