package auth

import (
	"context"
	"strconv"
	"time"

	"github.com/Maycon-Santos/go-snake-backend/cache"
	"github.com/Maycon-Santos/go-snake-backend/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenDetails struct {
	AccessToken string
	AccessUuid  uint64
	ExpiresAt   int64
}

func CompareHashAndPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GeneratePasswordHash(password string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(encrypted), err
}

func CreateToken(
	expiresIn time.Duration,
	secret string,
	accountID string,
) (*TokenDetails, error) {
	accessUuid, err := uuid.Generate()
	if err != nil {
		return nil, err
	}

	tokenDetails := &TokenDetails{}
	tokenDetails.ExpiresAt = time.Now().Add(expiresIn).Unix()
	tokenDetails.AccessUuid = *accessUuid

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"access_uuid": accessUuid,
			"account_id":  accountID,
			"expires_at":  tokenDetails.ExpiresAt,
		})

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	tokenDetails.AccessToken = tokenStr

	return tokenDetails, nil
}

func CreateAuth(ctx context.Context, cacheClient cache.Client, accountID string, tokenDetails *TokenDetails) error {
	at := time.Unix(tokenDetails.ExpiresAt, 0)
	now := time.Now()

	errAccess := cacheClient.Set(ctx, strconv.FormatUint(tokenDetails.AccessUuid, 10), accountID, at.Sub(now))
	if errAccess != nil {
		return errAccess
	}

	return nil
}
