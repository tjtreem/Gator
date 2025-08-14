package main

import (
	"fmt"
	"time"
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/tjtreem/gator/internal/database"
	
	
)


func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
	    return fmt.Errorf("missing required argument: time_between_reqs")
	}

	timeBetweenReqsStr := cmd.Args[0]
	timeBetweenReqs, err := time.ParseDuration(timeBetweenReqsStr)
	if err != nil {
	    return fmt.Errorf("invalid duration: %w", err)
	}
	
	fmt.Printf("Collecting feeds every %s\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	defer ticker.Stop()

	for {
	    if err := ScrapeFeeds(s, cmd); err != nil {
		fmt.Printf("Error scraping feeds: %v\n", err)
	    }
	    <-ticker.C
	}	

}


func ScrapeFeeds(s *State, cmd Command) error {
	ctx := context.Background()

	feed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
	   return fmt.Errorf("could not scrape feed: %w", err)
	}
	
	err = s.Db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
	    return fmt.Errorf("could not mark feed as fetched: %w", err)
	}

	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
	   return fmt.Errorf("unable to fetch feed: %w", err)
	}
	
	for _, item := range rssFeed.Channel.Item {
	    var publishedAt time.Time

	    if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
	    	publishedAt = t	    	
            } else {
		publishedAt = time.Now().UTC()
	    }
	
	    _, err := s.Db.CreatePost(context.Background(), database.CreatePostParams{
	    	ID:		uuid.New(),
		CreatedAt:	time.Now().UTC(),
		UpdatedAt:	time.Now().UTC(),
		FeedID:		feed.ID,
		Title:		item.Title,
		Url:		item.Link,
		Description:	sql.NullString{
		    String:	item.Description,
		    Valid:	true,
		},
		PublishedAt:	publishedAt,
	    })

	    if err != nil {
	        if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		continue
	        }	
	    	log.Printf("Couldn't create post: %v", err)
	    	continue
	    }
	}
	return nil

}

