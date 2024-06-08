package routes

import (
	"log"
	"net/http"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/julienschmidt/httprouter"
)

type pattern struct {
	Source string `json:"source"`
	Type   string `json:"type"`
}

type availableSkinsResponse struct {
	Colors  map[string]string  `json:"colors"`
	Pattern map[string]pattern `json:"patterns"`
}

func AvailableSkins(container container.Container) httprouter.Handle {
	var skinsRepository db.SkinsRepository

	err := container.Retrieve(&skinsRepository)
	if err != nil {
		log.Fatal(err)
	}

	return func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		availableSkins := availableSkinsResponse{
			Colors:  make(map[string]string),
			Pattern: make(map[string]pattern),
		}

		colors, err := skinsRepository.GetAllColors(request.Context())
		if err != nil {
			handleError(request.Context(), err)
		}

		for _, c := range colors {
			availableSkins.Colors[c.ID] = c.Color
		}

		patterns, err := skinsRepository.GetAllPatterns(request.Context())
		if err != nil {
			handleError(request.Context(), err)
		}

		for _, p := range patterns {
			availableSkins.Pattern[p.ID] = pattern{
				Source: p.Source,
				Type:   p.Type,
			}
		}

		response := responseConfig{
			Body: responseBody{
				Success: true,
				Result:  availableSkins,
			},
		}

		if makeResponse(request.Context(), writer, response); err != nil {
			handleError(request.Context(), err)
		}
	}
}
