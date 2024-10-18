package cli

import (
	"context"

	"github.com/meraiku/aggregator/internal/app"
	"github.com/meraiku/aggregator/internal/repository/sql"
)

type LoggedFunc func(state *app.State, cmd Command, user sql.GetUserRow) error

func MiddlewreLoggedIn(handler LoggedFunc) CLIFunc {
	return func(s *app.State, c Command) error {

		ctx := context.Background()

		user, err := s.Db.GetUser(ctx, s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, c, user)
	}
}
