package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/4madness7/gator/internal/config"
	"github.com/4madness7/gator/internal/database"
	_ "github.com/lib/pq"
)

const (
	dbUrl = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	st := state{
		db:  database.New(db),
		cfg: &cfg,
	}

	cmds := commands{
		funcs: map[string]func(*state, command) error{},
	}
	cmds.register("login", loginHandler)
	cmds.register("register", registerHander)
	cmds.register("reset", resetHandler)
	cmds.register("users", usersHandler)
	cmds.register("agg", aggHandler)
	cmds.register("addfeed", middlewareLoggedIn(addfeedHander))
	cmds.register("feeds", feedsHandler)
	cmds.register("follow", middlewareLoggedIn(followHandler))
	cmds.register("following", middlewareLoggedIn(followingHandler))
	cmds.register("unfollow", middlewareLoggedIn(unfollowHandler))

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
