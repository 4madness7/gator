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
	err := s.cfg.SetUser(username)
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
