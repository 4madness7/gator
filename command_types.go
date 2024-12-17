package main

import (
	"errors"
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	funcs map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
    c.funcs[name] = f
}


func (c *commands) run(s *state, cmd command) error {
    fun, ok := c.funcs[cmd.name]
    if !ok {
        return errors.New(fmt.Sprintf("'%s' command does not exist.", cmd.name))
    }
    err := fun(s, cmd)
    if err != nil {
        return err
    }
    return nil
}
