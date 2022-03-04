package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Maycon-Santos/go-snake-backend/cache"
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/Maycon-Santos/go-snake-backend/process"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
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
	ctrl := gomock.NewController(t)

	env, err := process.NewEnv()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when retrieve env vars.", err)
	}

	// dbConn, mock, err := sqlmock.New()
	// if err != nil {
	// 	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	// }

	// defer dbConn.Close()

	// cacheClient := cache.NewMockClient(ctrl)
	// accountsRepository := db.NewAccountsRepository(dbConn)
	// env, err := process.NewEnv()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// query := "SELECT id, username, password FROM accounts WHERE username=$1"

	// rows := sqlmock.NewRows([]string{"id", "username", "password"}).
	// 	AddRow("mocked_id", "mocked_username", "mocked_password")

	// mock.ExpectQuery(query).WithArgs("mocked_username").WillReturnRows(rows)

	// dependenciesContainer.Inject(env, &cacheClient, &accountsRepository)

	// signInHandler := SignInHandler(dependenciesContainer)
	// body := bytes.NewBufferString("")

	// resp, _ := doRequest("POST", "/v1/signin", body, signInHandler)

	// fmt.Print(resp)

	t.Run("should return an error when payload is invalid", func(t *testing.T) {
		dbConn, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection.", err)
		}

		defer dbConn.Close()

		dependenciesContainer := container.New()
		cacheClient := cache.NewMockClient(ctrl)
		accountsRepository := db.NewAccountsRepository(dbConn)

		dependenciesContainer.Inject(env, &cacheClient, &accountsRepository)

		signInHandler := SignInHandler(dependenciesContainer)
		requestBody := bytes.NewBufferString("invalid_payload")

		response, err := doRequest("POST", "/v1/signin", requestBody, signInHandler)
		if err != nil {
			t.Fatalf("%s: an error '%s' was not expected.", t.Name(), err)
		}

		var responseBody responseSchema

		expectedResponseBody := responseSchema{
			Success: false,
			Type:    TYPE_PAYLOAD_INVALID,
			Message: "playload invalid",
		}

		err = json.Unmarshal(response.Body.Bytes(), &responseBody)
		if err != nil {
			t.Fatalf("%s: an error '%s' was not expected.", t.Name(), err)
		}

		assert.Equal(t, expectedResponseBody, responseBody)
		assert.Equal(t, response.Result().StatusCode, http.StatusUnprocessableEntity)
	})

	t.Run("should return an error when username is missing", func(t *testing.T) {
		dbConn, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection.", err)
		}

		defer dbConn.Close()

		dependenciesContainer := container.New()
		cacheClient := cache.NewMockClient(ctrl)
		accountsRepository := db.NewAccountsRepository(dbConn)

		dependenciesContainer.Inject(env, &cacheClient, &accountsRepository)

		signInHandler := SignInHandler(dependenciesContainer)

		requestBody, err := json.Marshal(signInRequestBody{
			Username: "",
			Password: "",
		})
		if err != nil {
			t.Fatalf("%s: an error '%s' was not expected.", t.Name(), err)
		}

		response, err := doRequest("POST", "/v1/signin", bytes.NewBuffer(requestBody), signInHandler)
		if err != nil {
			t.Fatalf("%s: an error '%s' was not expected.", t.Name(), err)
		}

		var responseBody responseSchema

		expectedResponseBody := responseSchema{
			Success: false,
			Type:    TYPE_MISSING_USERNAME,
			Message: err.Error(),
		}

		err = json.Unmarshal(response.Body.Bytes(), &responseBody)
		if err != nil {
			t.Fatalf("%s: an error '%s' was not expected.", t.Name(), err)
		}

		assert.Equal(t, expectedResponseBody, responseBody)
		assert.Equal(t, response.Result().StatusCode, http.StatusForbidden)
	})

	// t.Run("should return an error when payload is invalid", func(t *testing.T) {
	// 	dbConn, mock, err := sqlmock.New()
	// 	if err != nil {
	// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	// 	}

	// 	defer dbConn.Close()

	// 	cacheClient := cache.NewMockClient(ctrl)
	// 	accountsRepository := db.NewAccountsRepository(dbConn)
	// 	env, err := process.NewEnv()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	mock.
	// 		ExpectQuery("SELECT id, username, password FROM accounts WHERE username=$1").
	// 		WithArgs("mocked_username").
	// 		WillReturnRows(
	// 			sqlmock.
	// 				NewRows([]string{"id", "username", "password"}).
	// 				AddRow("mocked_id", "mocked_username", "mocked_password"),
	// 		)

	// 	dependenciesContainer.Inject(env, &cacheClient, &accountsRepository)

	// 	signInHandler := SignInHandler(dependenciesContainer)
	// 	body := bytes.NewBufferString("")

	// 	resp, _ := doRequest("POST", "/v1/signin", body, signInHandler)

	// 	fmt.Print(resp)
	// })
}
