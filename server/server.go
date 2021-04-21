package server

import (
	"fmt"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/process"
)

func Listen(container container.Container) error {
	var env process.Env

	container.Retrieve(&env)

	routes := newRoutes(container)

	return http.ListenAndServe(fmt.Sprintf(":%d", env.ServerPort), routes)
}
