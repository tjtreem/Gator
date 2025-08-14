package main

import (
        "fmt"
        "context"
	"strconv"

        "github.com/tjtreem/gator/internal/database"
)



func HandlerBrowse(s *State, cmd Command, user database.User) error {
	limit := 2

	if len(cmd.Args) == 1 {
	    specifiedLimit, err := strconv.Atoi(cmd.Args[0])
	    if err != nil {
		return fmt.Errorf("Couldn't convert to integer: %w", err)
	    }

	    limit = specifiedLimit
	}
	
	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID:	user.ID,
		Limit:	int32(limit),
	})

	if err != nil {
	    return fmt.Errorf("Couldn't get posts: %w", err)
	}

	fmt.Printf("found %d posts for user %s:\n", len(posts), user.Name)

	for _, post := range posts {
	    fmt.Printf("Title: %s\n", post.Title)
	    fmt.Printf("URL: %s\n", post.Url)
	    fmt.Printf("Description: %s\n", post.Description.String)
	    fmt.Printf("Published at: %s\n", post.PublishedAt.Format("2006-01-02"))
	    fmt.Println("---")
	}
	return nil

}













