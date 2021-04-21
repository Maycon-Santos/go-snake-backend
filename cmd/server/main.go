package main

import (
	"context"
	"log"

	"github.com/Maycon-Santos/go-snake-backend/cache"
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/Maycon-Santos/go-snake-backend/server"
)

func main() {
	env, err := process.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.NewConnection(*env)
	if err != nil {
		log.Fatal(err)
	}

	accountsRepository := db.NewAccountsRepository(dbConn)
	cache, err := cache.NewClient(context.Background(), env.RedisAddress)
	if err != nil {
		log.Fatal(err)
	}

	dependenciesContainer := container.New()

	dependenciesContainer.Inject(
		env,
		cache,
		accountsRepository,
	)

	err = server.Listen(dependenciesContainer)
	if err != nil {
		log.Fatal(err)
	}
}
