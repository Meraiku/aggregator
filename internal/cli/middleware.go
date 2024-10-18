package cli

import (
	"context"

	"github.com/meraiku/aggregator/internal/repository/sql"
)

type LoggedFunc func(state *State, cmd Command, user sql.User) error

func MiddlewreLoggedIn(handler LoggedFunc) CLIFunc {
	return func(s *State, c Command) error {

		ctx := context.Background()

		user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, c, user)
	}
}
