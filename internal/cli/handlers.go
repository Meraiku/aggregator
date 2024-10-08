package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/aggregator/internal/repository/sql"
	"github.com/meraiku/aggregator/internal/rss"
)

func Login(state *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return ErrNoArgs
	}
	ctx := context.Background()

	user, err := state.db.GetUser(ctx, cmd.Args[0])
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return ErrUserNotExists
		}
		return err
	}

	err = state.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' has been set!\n", state.cfg.CurrentUserName)

	return nil
}

func Register(state *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return ErrNoArgs
	}

	ctx := context.Background()

	user := sql.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	userSQL, err := state.db.CreateUser(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrUserAlreadyExists
		}
		return err
	}

	if err := state.cfg.SetUser(userSQL.Name); err != nil {
		return err
	}

	fmt.Printf("User '%s' has been created!\n", userSQL.Name)

	return nil
}

func Reset(state *State, cmd Command) error {

	ctx := context.Background()

	if err := state.db.ResetUsers(ctx); err != nil {
		return err
	}

	fmt.Printf("Users succsessfuly has been deleted!\n")

	return nil
}

func Users(state *State, cmd Command) error {

	ctx := context.Background()

	users, err := state.db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user == state.cfg.CurrentUserName {
			user = fmt.Sprintf("%s (current)", user)
		}
		fmt.Printf("* %s\n", user)
	}

	return nil
}

func Agg(state *State, cmd Command) error {

	ctx := context.Background()

	feed, err := rss.FetchRSS(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func AddFeed(state *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return ErrInvalidArgumentsCount
	}

	ctx := context.Background()

	user, err := state.db.GetUser(ctx, state.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	params := sql.CreateFeedParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feed, err := state.db.CreateFeed(ctx, params)
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func Feeds(state *State, cmd Command) error {

	ctx := context.Background()

	data, err := state.db.GetAllFeeds(ctx)
	if err != nil {
		return err
	}

	for i := range data {
		fmt.Printf("%s\n%s\n%s\n", data[i].Name, data[i].Url, data[i].Name_2)
	}

	return nil
}
