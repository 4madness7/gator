package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/4madness7/gator/internal/database"
	"github.com/google/uuid"
)

func loginHandler(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("'login' expects 1 <username> parameter.")
	}

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User %s does not exist in the database.\nError: %w", username, err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User as been set to '%s'\n", username)
	return nil
}

func registerHander(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("'register' expects 1 <username> parameter.")
	}
	current_time := time.Now()
	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: current_time,
		UpdatedAt: current_time,
		Name:      cmd.args[0],
	}

	newUser, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(newUser.Name)
	if err != nil {
		return err
	}

	fmt.Println("User created and logged successfully.")
	fmt.Println("===== DEBUG =====")
	fmt.Printf(`User Data {
    id:         %v
    created at: %v
    updated at: %v
    name:       %s
}`,
		newUser.ID, newUser.CreatedAt, newUser.UpdatedAt, newUser.Name)
	fmt.Println()

	return nil
}

func resetHandler(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("'reset' does not expect any arguments.")
	}
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Could not reset DB: %w", err)
	}

	fmt.Println("Datadase reset successfully.")
	return nil
}

func usersHandler(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("'users' does not expect any arguments.")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Something went wrong: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("No users in the database")
	}

	for _, user := range users {
		currentStr := ""
		if s.cfg.CurrentUserName == user.Name {
			currentStr = " (current)"
		}
		fmt.Printf("* %s%s\n", user.Name, currentStr)
	}
	return nil
}

func aggHandler(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("'agg' does not expect any arguments.")
	}
    feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
    if err != nil {
        return err
    }
    fmt.Println(*feed)
	return nil
}
