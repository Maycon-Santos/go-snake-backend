package routes

import (
	"context"
	"encoding/json"
	"net/http"
)

type responseHeader struct {
	Status int
}

type responseBody struct {
	Success bool         `json:"success"`
	Type    responseType `json:"type,omitempty"`
	Message string       `json:"message,omitempty"`
	Result  interface{}  `json:"result,omitempty"`
}

type responseConfig struct {
	Header responseHeader
	Body   responseBody
}

type responseType string

const (
	TYPE_UNKNOWN                 = responseType("UNKNOWN")
	TYPE_ACCOUNT_NOT_FOUND       = responseType("ACCOUNT_NOT_FOUND")
	TYPE_ACCOUNT_PASSWORD_WRONG  = responseType("ACCOUNT_PASSWORD_WRONG")
	TYPE_ACCOUNT_USERNAME_EXISTS = responseType("ACCOUNT_USERNAME_EXISTS")
	TYPE_PAYLOAD_INVALID         = responseType("PAYLOAD_INVALID")
	TYPE_USERNAME_MISSING        = responseType("USERNAME_MISSING")
	TYPE_USERNAME_BELOW_MIN_LEN  = responseType("USERNAME_BELOW_MIN_LEN")
	TYPE_USERNAME_ABOVE_MAX_LEN  = responseType("USERNAME_ABOVE_MAX_LEN")
	TYPE_USERNAME_INVALID_CHAR   = responseType("USERNAME_INVALID_CHAR")
	TYPE_PASSWORD_MISSING        = responseType("PASSWORD_MISSING")
	TYPE_PASSWORD_BELOW_MIN_LEN  = responseType("PASSWORD_BELOW_MIN_LEN")
	TYPE_PASSWORD_ABOVE_MAX_LEN  = responseType("PASSWORD_ABOVE_MAX_LEN")
	TYPE_ROOM_NOT_FOUND          = responseType("TYPE_ROOM_NOT_FOUND")
)

func makeResponse(ctx context.Context, writer http.ResponseWriter, response responseConfig) error {
	resp, err := json.Marshal(response.Body)
	if err != nil {
		return err
	}

	status := http.StatusOK
	if response.Header.Status != 0 {
		status = response.Header.Status
	}

	writer.WriteHeader(status)

	_, err = writer.Write(resp)
	if err != nil {
		return err
	}

	return nil
}
