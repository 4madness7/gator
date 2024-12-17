package main

import (
	"fmt"
	"os"

	"github.com/4madness7/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	st := state{
		cfg: &cfg,
	}

	cmds := commands{
		funcs: map[string]func(*state, command) error{},
	}
	cmds.register("login", loginHandler)

	if len(os.Args) < 2 {
		fmt.Println("Please provide an argument.")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
