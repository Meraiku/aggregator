package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/aggregator/internal/app"
	"github.com/meraiku/aggregator/internal/repository/sql"
	"github.com/meraiku/aggregator/internal/rss"
)

func Login(state *app.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return ErrNoArgs
	}
	ctx := context.Background()

	user, err := state.Db.GetUser(ctx, cmd.Args[0])
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return ErrUserNotExists
		}
		return err
	}

	err = state.Cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' has been set!\n", state.Cfg.CurrentUserName)

	return nil
}

func Register(state *app.State, cmd Command) error {
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

	userSQL, err := state.Db.CreateUser(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrUserAlreadyExists
		}
		return err
	}

	if err := state.Cfg.SetUser(userSQL.Name); err != nil {
		return err
	}

	fmt.Printf("User '%s' has been created!\n", userSQL.Name)

	return nil
}

func Reset(state *app.State, cmd Command) error {

	ctx := context.Background()

	if err := state.Db.ResetUsers(ctx); err != nil {
		return err
	}

	fmt.Printf("Users succsessfuly has been resetted!\n")

	return nil
}

func Users(state *app.State, cmd Command) error {

	ctx := context.Background()

	users, err := state.Db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user == state.Cfg.CurrentUserName {
			user = fmt.Sprintf("%s (current)", user)
		}
		fmt.Printf("* %s\n", user)
	}

	return nil
}

func Agg(state *app.State, cmd Command) error {

	if len(cmd.Args) == 0 {
		return ErrNoArgs
	}

	dur, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	if dur < time.Minute {
		return fmt.Errorf("DDOS WARNING!!! Minimum duration is 1 minute. You entered %v\n", dur.String())
	}

	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	fmt.Printf("Scrapping feeds every %v\n", dur.String())

	for ; ; <-ticker.C {
		rss.ScrapeFeeds(state)
	}

}

func AddFeed(state *app.State, cmd Command, user sql.GetUserRow) error {
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

	feed, err := state.Db.CreateFeed(ctx, params)
	if err != nil {
		return err
	}

	_, err = state.Db.CreateFeedFollow(ctx, sql.CreateFeedFollowParams{
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

func Feeds(state *app.State, cmd Command) error {

	ctx := context.Background()

	data, err := state.Db.GetAllFeeds(ctx)
	if err != nil {
		return err
	}

	fmt.Println("----------------")
	for i := range data {
		fmt.Printf("Feed: %s\nURL: %s\nCreated by: %s\n---------------\n", data[i].FeedName, data[i].Url, data[i].UserName)
	}

	return nil
}

func Follow(state *app.State, cmd Command, user sql.GetUserRow) error {

	if len(cmd.Args) == 0 {
		return ErrNoArgs
	}

	ctx := context.Background()

	feedID, err := state.Db.GetFeedIDByURL(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	data, err := state.Db.CreateFeedFollow(ctx, sql.CreateFeedFollowParams{
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

func Following(state *app.State, cmd Command, user sql.GetUserRow) error {

	ctx := context.Background()

	data, err := state.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for i := range data {
		fmt.Printf("Feed '%s' followed by '%s'\n", data[i].FeedName, data[i].UserName)
	}

	return nil
}

func Unfollow(state *app.State, cmd Command, user sql.GetUserRow) error {

	if len(cmd.Args) == 0 {
		return ErrNoArgs
	}

	ctx := context.Background()

	if err := state.Db.DeleteFeedFollow(ctx, sql.DeleteFeedFollowParams{
		Url:    cmd.Args[0],
		UserID: user.ID,
	}); err != nil {
		return err
	}

	fmt.Printf("Feed '%s' has been unfollowed!\n", cmd.Args[0])

	return nil
}

func Browse(state *app.State, cmd Command, user sql.GetUserRow) error {

	limit, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		limit = 2
	}

	if len(cmd.Args) == 0 {
		limit = 2
	}

	ctx := context.Background()

	posts, err := state.Db.GetPostsForUser(ctx, sql.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println("----------------")
		fmt.Printf("Title: %s\n\nDescription: %s\nURL: %s\n--------------\n", post.Title, post.Description.String, post.Url)
	}
	return nil
}
