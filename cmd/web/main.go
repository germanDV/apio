package main

import (
	"github.com/a-h/rest"

	"github.com/germandv/apio/internal/cache/memorycache"
	"github.com/germandv/apio/internal/config"
	"github.com/germandv/apio/internal/logger"
	"github.com/germandv/apio/internal/memorydb"
	"github.com/germandv/apio/internal/notes"
	"github.com/germandv/apio/internal/tags"
	"github.com/germandv/apio/internal/tokenauth"
	"github.com/germandv/apio/internal/web"
)

func main() {
	cfg := config.Get()

	logger, err := logger.Get(cfg.LogFormat, cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	auth, err := tokenauth.New(cfg.AuthPrivKey, cfg.AuthPublKey)
	if err != nil {
		panic(err)
	}

	cacheClient, err := memorycache.New()
	if err != nil {
		panic(err)
	}

	// dbPool, err := db.InitWithConnStr(cfg.PostgresConnStr)
	// if err != nil {
	// 	panic(err)
	// }

	tagSvc := tags.NewService(memorydb.NewTagsRepository())
	noteSvc := notes.NewService(memorydb.NewNotesRepository())

	oas := rest.NewAPI("apio")
	api := web.New(logger, auth, oas, tagSvc, noteSvc, cacheClient)

	api.ListenAndServe()
}
