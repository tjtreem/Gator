package main

import (
	"os"
	"fmt"
	"time"
	"errors"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/tjtreem/gator/internal/config"
	"github.com/tjtreem/gator/internal/database"
	
)


type State struct {
	Db *database.Queries
	Cfg *config.Config
}


type Command struct {
	Name		string
	Args		[]string
}


type Commands struct {
	Handlers map[string]func(*State, Command) error
}




func (c *Commands) Run(s *State, cmd Command) error {
	handler, exists := c.Handlers[cmd.Name]
	if !exists {
	    return errors.New("command not found")
	}
	return handler(s, cmd)

}


func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Handlers[name] = f
}


func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
	    return errors.New("Login requires a username")
	}
	
	username := cmd.Args[0]
	ctx := context.Background()

	_, err := s.Db.GetUser(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
	    fmt.Println("User does not exist")
	    os.Exit(1)
	} else if err != nil {
	    return err
	}


	err = s.Cfg.SetUser(username)
	if err != nil {
	    return err
	}
	
	fmt.Println("User has been set")
	return nil
}


func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
	    return errors.New("register requires a username")
	}

	username := cmd.Args[0]
	ctx := context.Background()
	_, err := s.Db.GetUser(ctx, username)
	if err == nil {
	    fmt.Println("User already exists")
	    os.Exit(1)
	} else if !errors.Is(err, sql.ErrNoRows) {
	    return err
	}
	
	
	now := time.Now()
	params := database.CreateUserParams{
	    ID:		uuid.New(),
	    CreatedAt:	now,
	    UpdatedAt:	now,
	    Name:	username,
	}
	
	user, err := s.Db.CreateUser(ctx, params)
	if err != nil {
	    return err
	}

	err = s.Cfg.SetUser(username)
	if err != nil {
	    return err
	}

	fmt.Println("User created:", username)
	fmt.Printf("User data: %+v\n", user)

	return nil

}



func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
	    return fmt.Errorf("couldn't delete users: %w", err)
	}
	fmt.Println("Database reset successfully!")
	return nil
}


func HandlerGetUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
	    return fmt.Errorf("couldn't get users: %w", err)
	}

	cfg, err := config.Read()
	if err != nil {
	    return fmt.Errorf("couldn't read config: %w", err)
	}

	for _, user := range users {
	    if user.Name == cfg.CurrentUserName {
		fmt.Printf("* %s (current)\n", user.Name)
	    } else {
		fmt.Printf("* %s\n", user.Name)
	    }
	}
	return nil
}


func HandlerAgg(s *State, cmd Command) error {
	ctx := context.Background()
	agg, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
	    return fmt.Errorf("unable to fetch URL: %s", err)
	}
	
	fmt.Printf("%s\n\n", agg.Channel.Title)
	fmt.Printf("%s\n\n", agg.Channel.Description)

	for _, item := range agg.Channel.Item {
	    fmt.Printf("TITLE: %s\n", item.Title)
	    fmt.Printf("LINK: %s\n", item.Link)
	    fmt.Printf("DATE: %s\n", item.PubDate)
	    fmt.Printf("DESCRIPTION: %s\n\n", item.Description)
	}

	return nil
}
