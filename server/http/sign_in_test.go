package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Maycon-Santos/go-snake-backend/cache"
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
)

func doRequest(method, uri string, body *bytes.Buffer, handle httprouter.Handle) (*httptest.ResponseRecorder, error) {
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

func TestSignInHandler(t *testing.T) {
	dependenciesContainer := container.New()
	ctrl := gomock.NewController(t)
	cacheClient := cache.NewMockClient(ctrl)
	env := process.Env{
		JWT: process.JWT{
			ExpiresIn:        time.Duration(time.Second * 60),
			RefreshExpiresIn: time.Duration(time.Second * 60),
			Secret:           "secret",
			RefreshSecret:    "refresh_secret",
		},
	}

	dependenciesContainer.Inject(&env, &cacheClient)

	signInHandler := SignInHandler(dependenciesContainer)
	body := bytes.NewBufferString("")

	resp, _ := doRequest("POST", "/v1/signin", body, signInHandler)

	fmt.Print(resp)
}
