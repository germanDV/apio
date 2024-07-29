package main

import (
	"fmt"
	"os"

	"github.com/germandv/apio/internal/config"
	"github.com/germandv/apio/internal/tokenauth"
)

// Generate auth token
func main() {
	if len(os.Args) < 3 {
		panic("provide a user ID and a role")
	}

	userID := os.Args[1]
	role := os.Args[2]

	cfg := config.Get()
	tokenService, err := tokenauth.New(cfg.AuthPrivKey, cfg.AuthPublKey)
	if err != nil {
		panic(err)
	}

	token, err := tokenService.Generate(userID, role)
	if err != nil {
		panic(err)
	}

	fmt.Println(token)
}
