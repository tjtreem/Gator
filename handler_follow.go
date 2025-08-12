package main

import (
	"fmt"
	"time"
	"context"
	"github.com/google/uuid"
	"github.com/tjtreem/gator/internal/database"
	
)

func HandlerFollow (s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
	    return fmt.Errorf("URL must be provided")
	}

	url := cmd.Args[0]
	ctx := context.Background()

	feed, err := s.Db.GetFeedByUrl(ctx, url)
	if err != nil {
	    return fmt.Errorf("No url found in the database")
	}
	
	
	id := uuid.New()
	now := time.Now()

	feedfollow, err := s.Db.CreateFeedFollow(
	    context.Background(),
	    database.CreateFeedFollowParams{
		ID:		id,
		CreatedAt:	now,
		UpdatedAt:	now,
		UserID:		user.ID,
		FeedID:		feed.ID,
	    },
	)

	if err != nil {
	    return fmt.Errorf("unable to insert follow record")
	}

	fmt.Printf("User %s is now following %s", feedfollow.UserName, feedfollow.FeedName)
	
	return nil

}


func HandlerFollowing (s *State, cmd Command, user database.User) error {
	ctx := context.Background()

	follows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
	    return fmt.Errorf("no user found in the database")
	}

	for _, follow := range follows {
	    fmt.Println(follow.FeedName)
	}
	
	return nil

}


func HandlerUnfollow (s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
	    return fmt.Errorf("URL must be provided")
	}

	url := cmd.Args[0]
	ctx := context.Background()

	feed, err := s.Db.GetFeedByUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("No url found in the database: %w", err)
	}
	
	
	err = s.Db.DeleteFeedFollow(
	    context.Background(),
	    database.DeleteFeedFollowParams{
		UserID:		user.ID,
		FeedID:		feed.ID,
	    },
	)

	if err != nil {
	    return fmt.Errorf("unable to unfollow record")
	}

	fmt.Printf("User %s is now unfollowing %s", user.Name, feed.Name)
	
	return nil

}


