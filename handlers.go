package main

import (
	"errors"
	"fmt"
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
