package rss

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/aggregator/internal/app"
	repo "github.com/meraiku/aggregator/internal/repository/sql"
)

func ScrapeFeeds(state *app.State) {

	ctx := context.Background()

	feed, err := state.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return
	}

	if err := state.Db.MarkFetched(ctx, repo.MarkFetchedParams{
		ID:            feed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:     time.Now(),
	}); err != nil {
		return
	}

	fmt.Println("------------------")
	fmt.Printf("Starting to fetch feed '%s'...\n", feed.Url)
	fmt.Println("------------------")

	rssFeed, err := FetchRSS(ctx, feed.Url)
	if err != nil {
		return
	}

	for _, item := range rssFeed.Channel.Item {

		pubDate := sql.NullTime{}
		pd, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err == nil {
			pubDate.Time = pd
			pubDate.Valid = true
		}

		err = state.Db.CreatePost(ctx, repo.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: pubDate,
			FeedID:      feed.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Println(err)
		}
	}

	fmt.Printf("Posts from feed '%s' has been fetched!\n", feed.Url)
}
