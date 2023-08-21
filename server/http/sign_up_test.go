package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Maycon-Santos/go-snake-backend/cache"
	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/Maycon-Santos/go-snake-backend/db"
	"github.com/Maycon-Santos/go-snake-backend/process"
	test_utils "github.com/Maycon-Santos/go-snake-backend/test_utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	dependenciesContainer := container.New()
	cacheClient := cache.NewMockClient(ctrl)
	accountsRepository := db.NewMockAccountsRepository(ctrl)
	env := &process.Env{
		JWT: process.JWT{
			ExpiresIn:        time.Duration(time.Second * 60),
			RefreshExpiresIn: time.Duration(time.Second * 60),
			Secret:           "secret",
			RefreshSecret:    "refresh_secret",
		},
	}

	dependenciesContainer.Inject(&accountsRepository, &cacheClient, env)

	signUpHandler := SignUpHandler(dependenciesContainer)

	t.Run("should response the tokens with success", func(t *testing.T) {
		reqBody, _ := json.Marshal(signInRequestBody{
			Username: "michael",
			Password: "123456",
		})

		accountsRepository.EXPECT().CheckUsernameExists(gomock.Any(), gomock.Eq("michael")).Return(false, nil)

		accountsRepository.EXPECT().Save(gomock.Any(), gomock.Eq("michael"), gomock.Any()).Return("8", nil)

		cacheClient.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)

		resRecorder, _ := test_utils.DoRequest("POST", "/v1/signin", bytes.NewBuffer(reqBody), signUpHandler)

		var resBody responseBody
		json.Unmarshal(resRecorder.Body.Bytes(), &resBody)

		result := resBody.Result.(map[string]interface{})

		assert.Equal(t, http.StatusCreated, resRecorder.Result().StatusCode)
		assert.True(t, resBody.Success)
		assert.NotEmpty(t, result["access_token"])
		assert.NotEmpty(t, result["refresh_token"])
	})

	t.Run("should response an existing user message", func(t *testing.T) {
		reqBody, _ := json.Marshal(signInRequestBody{
			Username: "michael",
			Password: "123456",
		})

		accountsRepository.EXPECT().CheckUsernameExists(gomock.Any(), gomock.Eq("michael")).Return(true, nil)

		resRecorder, _ := test_utils.DoRequest("POST", "/v1/signin", bytes.NewBuffer(reqBody), signUpHandler)

		var resBody responseBody
		json.Unmarshal(resRecorder.Body.Bytes(), &resBody)

		assert.Equal(t, http.StatusUnauthorized, resRecorder.Result().StatusCode)

		assert.Equal(
			t,
			responseBody{
				Success: false,
				Type:    TYPE_ACCOUNT_USERNAME_EXISTS,
				Message: "username already in use",
			},
			resBody,
		)
	})

	t.Run("should response an invalid payload message", func(t *testing.T) {
		resRecorder, _ := test_utils.DoRequest("POST", "/v1/signin", bytes.NewBufferString("{invalid}"), signUpHandler)
		var resBody responseBody

		json.Unmarshal(resRecorder.Body.Bytes(), &resBody)

		assert.Equal(t, http.StatusUnprocessableEntity, resRecorder.Result().StatusCode)

		assert.Equal(
			t,
			responseBody{
				Success: false,
				Type:    TYPE_PAYLOAD_INVALID,
				Message: "playload is invalid",
			},
			resBody,
		)
	})

	t.Run("should response an error when `username` field is less than 4", func(t *testing.T) {
		reqBody, _ := json.Marshal(signInRequestBody{
			Username: "u",
			Password: "password",
		})

		resRecorder, _ := test_utils.DoRequest("POST", "/v1/signin", bytes.NewBuffer(reqBody), signUpHandler)
		var resBody responseBody

		json.Unmarshal(resRecorder.Body.Bytes(), &resBody)

		assert.Equal(t, http.StatusForbidden, resRecorder.Result().StatusCode)

		assert.Equal(
			t,
			responseBody{
				Success: false,
				Type:    TYPE_USERNAME_BELOW_MIN_LEN,
				Message: "the username field must be greater than 4",
			},
			resBody,
		)
	})

	t.Run("should response an error when `username` field is greater than 15", func(t *testing.T) {
		reqBody, _ := json.Marshal(signInRequestBody{
			Username: "username_123456789",
			Password: "password",
		})

		resRecorder, _ := test_utils.DoRequest("POST", "/v1/signin", bytes.NewBuffer(reqBody), signUpHandler)
		var resBody responseBody

		json.Unmarshal(resRecorder.Body.Bytes(), &resBody)

		assert.Equal(t, http.StatusForbidden, resRecorder.Result().StatusCode)

		assert.Equal(
			t,
			responseBody{
				Success: false,
				Type:    TYPE_USERNAME_ABOVE_MAX_LEN,
				Message: "the username field must be less than 15",
			},
			resBody,
		)
	})

	t.Run("should response an error when `username` field has an invalid char", func(t *testing.T) {
		reqBody, _ := json.Marshal(signInRequestBody{
			Username: "user name",
			Password: "password",
		})

		resRecorder, _ := test_utils.DoRequest("POST", "/v1/signin", bytes.NewBuffer(reqBody), signUpHandler)
		var resBody responseBody

		json.Unmarshal(resRecorder.Body.Bytes(), &resBody)

		assert.Equal(t, http.StatusForbidden, resRecorder.Result().StatusCode)

		assert.Equal(
			t,
			responseBody{
				Success: false,
				Type:    TYPE_USERNAME_INVALID_CHAR,
				Message: "username field cannot have the following characters:  ",
			},
			resBody,
		)
	})

	t.Run("should response an error when `username` field is missing", func(t *testing.T) {
		reqBody, _ := json.Marshal(signInRequestBody{
			Username: "",
			Password: "password",
		})

		resRecorder, _ := test_utils.DoRequest("POST", "/v1/signin", bytes.NewBuffer(reqBody), signUpHandler)
		var resBody responseBody

		json.Unmarshal(resRecorder.Body.Bytes(), &resBody)

		assert.Equal(t, http.StatusForbidden, resRecorder.Result().StatusCode)

		assert.Equal(
			t,
			responseBody{
				Success: false,
				Type:    TYPE_USERNAME_MISSING,
				Message: "missing username field",
			},
			resBody,
		)
	})

	t.Run("should response an error when `password` field is less than 6", func(t *testing.T) {
		reqBody, _ := json.Marshal(signInRequestBody{
			Username: "username",
			Password: "pass",
		})

		resRecorder, _ := test_utils.DoRequest("POST", "/v1/signin", bytes.NewBuffer(reqBody), signUpHandler)
		var resBody responseBody

		json.Unmarshal(resRecorder.Body.Bytes(), &resBody)

		assert.Equal(t, http.StatusForbidden, resRecorder.Result().StatusCode)

		assert.Equal(
			t,
			responseBody{
				Success: false,
				Type:    TYPE_PASSWORD_BELOW_MIN_LEN,
				Message: "the password field must be greater than 6",
			},
			resBody,
		)
	})

	t.Run("should response an error when `password` field is greater than 25", func(t *testing.T) {
		reqBody, _ := json.Marshal(signInRequestBody{
			Username: "username",
			Password: "password_123456789_123456789",
		})

		resRecorder, _ := test_utils.DoRequest("POST", "/v1/signin", bytes.NewBuffer(reqBody), signUpHandler)
		var resBody responseBody

		json.Unmarshal(resRecorder.Body.Bytes(), &resBody)

		assert.Equal(t, http.StatusForbidden, resRecorder.Result().StatusCode)

		assert.Equal(
			t,
			responseBody{
				Success: false,
				Type:    TYPE_PASSWORD_ABOVE_MAX_LEN,
				Message: "the password field must be less than 25",
			},
			resBody,
		)
	})

	t.Run("should response an error when `password` field is missing", func(t *testing.T) {
		reqBody, _ := json.Marshal(signInRequestBody{
			Username: "username",
			Password: "",
		})

		resRecorder, _ := test_utils.DoRequest("POST", "/v1/signin", bytes.NewBuffer(reqBody), signUpHandler)
		var resBody responseBody

		json.Unmarshal(resRecorder.Body.Bytes(), &resBody)

		assert.Equal(t, http.StatusForbidden, resRecorder.Result().StatusCode)

		assert.Equal(
			t,
			responseBody{
				Success: false,
				Type:    TYPE_PASSWORD_MISSING,
				Message: "missing password field",
			},
			resBody,
		)
	})
}
