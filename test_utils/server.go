package testutils

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/julienschmidt/httprouter"
)

func DoRequest(method, uri string, body *bytes.Buffer, handle httprouter.Handle) (*httptest.ResponseRecorder, error) {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	router := httprouter.New()
	router.Handle(method, uri, handle)
	router.ServeHTTP(resp, req)
	return resp, nil
}
