package auth

import (
	"context"
	"strconv"
	"time"

	"github.com/Maycon-Santos/go-snake-backend/cache"
	"github.com/Maycon-Santos/go-snake-backend/uuid"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   uint64
	RefreshUuid  uint64
	AtExpires    int64
	RtExpires    int64
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
	refreshExpiresIn time.Duration,
	secret string,
	refreshSecret string,
	accountID string,
) (*TokenDetails, error) {
	var err error

	accessUuid, err := uuid.Generate()
	if err != nil {
		return nil, err
	}

	refreshUuid, err := uuid.Generate()
	if err != nil {
		return nil, err
	}

	tokenDetails := &TokenDetails{}
	tokenDetails.AtExpires = time.Now().Add(expiresIn).Unix()
	tokenDetails.AccessUuid = *accessUuid

	tokenDetails.RtExpires = time.Now().Add(refreshExpiresIn).Unix()
	tokenDetails.RefreshUuid = *refreshUuid

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = tokenDetails.AccessUuid
	atClaims["account_id"] = accountID
	atClaims["exp"] = tokenDetails.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	if tokenDetails.AccessToken, err = at.SignedString([]byte(secret)); err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = tokenDetails.AccessUuid
	rtClaims["account_id"] = accountID
	rtClaims["exp"] = tokenDetails.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	tokenDetails.RefreshToken, err = rt.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}

	return tokenDetails, nil
}

func CreateAuth(ctx context.Context, cacheClient cache.Client, accountID string, tokenDetails *TokenDetails) error {
	at := time.Unix(tokenDetails.AtExpires, 0)
	rt := time.Unix(tokenDetails.RtExpires, 0)
	now := time.Now()

	errAccess := cacheClient.Set(ctx, strconv.FormatUint(tokenDetails.AccessUuid, 10), accountID, at.Sub(now))
	if errAccess != nil {
		return errAccess
	}

	errRefresh := cacheClient.Set(ctx, strconv.FormatUint(tokenDetails.RefreshUuid, 10), accountID, rt.Sub(now))
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}
