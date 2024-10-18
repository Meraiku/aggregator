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

	fmt.Printf("Users succsessfuly has been resetted!\n")

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

func AddFeed(state *State, cmd Command, user sql.User) error {
	if len(cmd.Args) < 2 {
		return ErrInvalidArgumentsCount
	}

	ctx := context.Background()

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

	_, err = state.db.CreateFeedFollow(ctx, sql.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed '%s' has been created!\n", feed.Name)

	return nil
}

func Feeds(state *State, cmd Command) error {

	ctx := context.Background()

	data, err := state.db.GetAllFeeds(ctx)
	if err != nil {
		return err
	}

	fmt.Println("----------------")
	for i := range data {
		fmt.Printf("Feed: %s\nURL: %s\nCreated by: %s\n---------------\n", data[i].FeedName, data[i].Url, data[i].UserName)
	}

	return nil
}

func Follow(state *State, cmd Command, user sql.User) error {

	ctx := context.Background()

	feedID, err := state.db.GetFeedIDByURL(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	data, err := state.db.CreateFeedFollow(ctx, sql.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feedID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed '%s' has been followed by '%s'!\n", data.FeedName, data.UserName)

	return nil
}

func Following(state *State, cmd Command, user sql.User) error {

	ctx := context.Background()

	data, err := state.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for i := range data {
		fmt.Printf("Feed '%s' followed by '%s'\n", data[i].FeedName, data[i].UserName)
	}

	return nil
}

func Unfollow(state *State, cmd Command, user sql.User) error {

	ctx := context.Background()

	if err := state.db.DeleteFeedFollow(ctx, sql.DeleteFeedFollowParams{
		Url:    cmd.Args[0],
		UserID: user.ID,
	}); err != nil {
		return err
	}

	fmt.Printf("Feed '%s' has been unfollowed!\n", cmd.Args[0])

	return nil
}
