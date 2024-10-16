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

	if userID == "" {
		fmt.Println("using hardcoded user ID")
		userID = "01929796-284c-7cab-9915-512e94297234"
	}

	if role == "" {
		fmt.Println("using `user` role")
		role = "user"
	}

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
