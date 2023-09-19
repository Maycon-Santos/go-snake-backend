package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/cache"
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/julienschmidt/httprouter"
)

type getUserResponseResult struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func GetUserHandler(container container.Container) httprouter.Handle {
	var (
		env                process.Env
		cache              cache.Client
		accountsRepository db.AccountsRepository
	)

	err := container.Retrieve(&env, &cache, &accountsRepository)
	if err != nil {
		log.Fatal(err)
	}

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("account_id")
		username := params.ByName("account_username")

		resBytes, err := json.Marshal(getUserResponseResult{
			ID:       id,
			Username: username,
		})
		if err != nil {
			handleError(request.Context(), err)
			return
		}

		writer.Write(resBytes)
	}
}
