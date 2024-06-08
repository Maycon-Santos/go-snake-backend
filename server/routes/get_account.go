package routes

import (
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/julienschmidt/httprouter"
)

type skinResult struct {
	Color   string `json:"color"`
	Pattern string `json:"pattern"`
}

type getUserResponseResult struct {
	ID       string     `json:"id"`
	Username string     `json:"username"`
	Skin     skinResult `json:"skin"`
}

func GetAccount(container container.Container) httprouter.Handle {
	var (
		accountsRepository db.AccountsRepository
		skinsRepository    db.SkinsRepository
	)

	err := container.Retrieve(&accountsRepository, &skinsRepository)
	if err != nil {
		log.Fatal(err)
	}

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		accountID := params.ByName("account_id")
		username := params.ByName("account_username")

		result := getUserResponseResult{
			ID:       accountID,
			Username: username,
		}

		skin, err := skinsRepository.GetAccountSkin(request.Context(), accountID)
		if err != nil {
			handleError(request.Context(), err)
			return
		}

		if skin != nil {
			result.Skin = skinResult{
				Color:   skin.ColorID,
				Pattern: skin.PatternID,
			}
		}

		response := responseConfig{
			Body: responseBody{
				Success: true,
				Result:  result,
			},
		}

		if err = makeResponse(request.Context(), writer, response); err != nil {
			handleError(request.Context(), err)
		}
	}
}
