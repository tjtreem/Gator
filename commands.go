package main

import (
	"errors"
	"context"
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


func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
    return func(s *State, cmd Command) error {
	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
	    return err
	}

	err = handler(s, cmd, user)
	if err != nil {
	   return err
	}
	return nil
    }
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











