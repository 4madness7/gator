package main

import (
	"fmt"
	"os"

	"github.com/4madness7/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

    err = cfg.SetUser("marco")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

    fmt.Println(cfg)
}
