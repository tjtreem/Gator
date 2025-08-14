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

