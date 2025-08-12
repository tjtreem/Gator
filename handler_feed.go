package main

import (
	"fmt"
	"time"
	"context"
	"github.com/google/uuid"
	"github.com/tjtreem/gator/internal/database"
)


func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
	    return fmt.Errorf("usage: addfeed <name> <url>")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	

	id := uuid.New()
	now := time.Now()

	feed, err := s.Db.CreateFeed(
	    context.Background(),
	    database.CreateFeedParams{
		ID:		id,
		CreatedAt:	now,
		UpdatedAt:	now,
		Name:		name,
		Url:		url,
		UserID:		user.ID,
	    },
	)
	if err != nil {
	    return fmt.Errorf("could not create feed: %w", err)
	}

	fmt.Printf("ID:		%s\n", feed.ID)
	fmt.Printf("Name:	%s\n", feed.Name)
	fmt.Printf("URL:	%s\n", feed.Url)
	fmt.Printf("User ID:	%s\n", feed.UserID)
	fmt.Printf("Created:	%v\n", feed.CreatedAt)
	fmt.Printf("Updated:	%v\n", feed.UpdatedAt)
	
	feedFollow, err := s.Db.CreateFeedFollow(
	    context.Background(),
	    database.CreateFeedFollowParams{
		ID:		uuid.New(),
		CreatedAt:	time.Now().UTC(),
		UpdatedAt:	time.Now().UTC(),
		UserID:		user.ID,
		FeedID:		feed.ID,
	    },
	)
	if err != nil {
	    return fmt.Errorf("couldn't create feed follow: %w", err)
	}
	
	fmt.Println("\nFeed followed successfully!")
	fmt.Printf("User: %s\n", feedFollow.UserName)
	fmt.Printf("Feed: %s\n", feedFollow.FeedName)


	return nil
}


func HandlerPrintFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
	    return fmt.Errorf("couldn't get feeds: %w", err)
	}

	if len(feeds) == 0 {
	    fmt.Println("No feeds found.")
	    return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
	    user, err := s.Db.GetUserById(context.Background(), feed.UserID)
	    if err != nil {
		return fmt.Errorf("couldn't get user: %w", err)
	    }

	    fmt.Printf("* Name: %s\n", feed.Name)
	    fmt.Printf("* URL: %s\n", feed.Url)
	    fmt.Printf("* User: %s\n", user.Name)
	    fmt.Println("=======================================")
	}

	return nil
}
	

