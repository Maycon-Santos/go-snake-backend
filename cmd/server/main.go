package main

import (
	"context"
	"log"

	"github.com/Maycon-Santos/go-snake-backend/cache"
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/Maycon-Santos/go-snake-backend/game"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/Maycon-Santos/go-snake-backend/server"
)

func main() {
	env, err := process.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.NewConnection(env)
	if err != nil {
		log.Fatal(err)
	}

	defer dbConn.Close()

	accountsRepository := db.NewAccountsRepository(dbConn)
	skinsRepository := db.NewSkinsRepository(dbConn)

	cacheClient, err := cache.NewClient(context.Background(), env.RedisAddress)
	if err != nil {
		log.Fatal(err)
	}

	dependenciesContainer := container.New()
	matches := game.NewMatches()

	err = dependenciesContainer.Inject(
		env,
		&cacheClient,
		&accountsRepository,
		&skinsRepository,
		&matches,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Listen(dependenciesContainer)
	if err != nil {
		log.Fatal(err)
	}
}
