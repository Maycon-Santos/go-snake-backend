package http

import "github.com/Maycon-Santos/go-snake-backend/validator"

var usernameValidator = validator.
	Field("username").
	Required().
	MinLen(4).
	MaxLen(15).
	NoContains([]string{" "})

var usernameResponseErrors = map[string]responseType{
	validator.Required:   TYPE_USERNAME_MISSING,
	validator.MinLen:     TYPE_USERNAME_BELOW_MIN_LEN,
	validator.MaxLen:     TYPE_USERNAME_ABOVE_MAX_LEN,
	validator.NoContains: TYPE_USERNAME_INVALID_CHAR,
}

var passwordValidator = validator.
	Field("password").
	Required().
	MinLen(6).
	MaxLen(25)

var passwordResponseErrors = map[string]responseType{
	validator.Required: TYPE_PASSWORD_MISSING,
	validator.MinLen:   TYPE_PASSWORD_BELOW_MIN_LEN,
	validator.MaxLen:   TYPE_PASSWORD_ABOVE_MAX_LEN,
}
