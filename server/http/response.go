package http

import (
	"context"
	"encoding/json"
	"net/http"
)

type responseConfig struct {
	Status  int          `json:"-"`
	Success bool         `json:"success"`
	Type    responseType `json:"type,omitempty"`
	Message string       `json:"message,omitempty"`
	Result  interface{}  `json:"result,omitempty"`
}

type responseType = string

const (
	TYPE_UNKNOWN                = responseType("UNKNOWN")
	TYPE_ACCOUNT_NOT_FOUND      = responseType("ACCOUNT_NOT_FOUND")
	TYPE_ACCOUNT_PASSWORD_WRONG = responseType("ACCOUNT_PASSWORD_WRONG")
	TYPE_USERNAME_EXISTS        = responseType("USERNAME_EXISTS")
	TYPE_PAYLOAD_INVALID        = responseType("PAYLOAD_INVALID")
	TYPE_MISSING_USERNAME       = responseType("MISSING_USERNAME")
	TYPE_BELOW_LEN_USERNAME     = responseType("BELOW_LEN_USERNAME")
	TYPE_ABOVE_LEN_USERNAME     = responseType("ABOVE_LEN_USERNAME")
	TYPE_MISSING_PASSWORD       = responseType("MISSING_PASSWORD")
	TYPE_BELOW_LEN_PASSWORD     = responseType("BELOW_LEN_PASSWORD")
	TYPE_ABOVE_LEN_PASSWORD     = responseType("ABOVE_LEN_PASSWORD")
)

func makeResponse(ctx context.Context, writer http.ResponseWriter, response responseConfig) error {
	resp, err := json.Marshal(response)
	if err != nil {
		return err
	}

	status := http.StatusOK
	if response.Status != 0 {
		status = response.Status
	}

	writer.WriteHeader(status)

	_, err = writer.Write(resp)
	if err != nil {
		return err
	}

	return nil
}
