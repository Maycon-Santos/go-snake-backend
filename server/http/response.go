package http

import (
	"context"
	"encoding/json"
	"net/http"
)

type responseConfig struct {
	Success bool         `json:"success"`
	Type    responseType `json:"type,omitempty"`
	Message string       `json:"message,omitempty"`
	Result  interface{}  `json:"result,omitempty"`
}

type responseType = string

const (
	TYPE_ACCOUNT_NOT_FOUND      = responseType("ACCOUNT_NOT_FOUND")
	TYPE_ACCOUNT_PASSWORD_WRONG = responseType("ACCOUNT_PASSWORD_WRONG")
	TYPE_USERNAME_EXISTS        = responseType("USERNAME_EXISTS")
)

func makeResponse(ctx context.Context, writer http.ResponseWriter, response responseConfig) error {
	resp, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = writer.Write(resp)
	if err != nil {
		return err
	}

	return nil
}
