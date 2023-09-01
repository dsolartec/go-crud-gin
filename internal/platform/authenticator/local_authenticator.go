package authenticator

import (
	"slices"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"go-crud-gin/internal/apperror"
	"go-crud-gin/internal/platform/logger"
)

const jwtSecret = "Abecedario"

type tokenClaims struct {
	jwt.StandardClaims
	Permissions []string `json:"permissions,omitempty"`
}

func (c tokenClaims) Validate() error {
	return nil
}

type localAuthenticator struct {
	logger logger.Logger
}

func (auth *localAuthenticator) GetToken(data AuthenticatorToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(data.UserID),
			ExpiresAt: time.Now().Add(time.Duration(1) * time.Hour).Unix(),
		},
		Permissions: data.Permissions,
	})

	return token.SignedString([]byte(jwtSecret))
}

func (auth *localAuthenticator) Authenticate(tokenStr string, permissions []string) (*AuthenticatorToken, error) {
	if !strings.HasPrefix(tokenStr, "Bearer") {
		return nil, apperror.NewErrUnauthorized()
	}

	var claims tokenClaims

	token, err := jwt.ParseWithClaims(strings.Replace(tokenStr, "Bearer ", "", 1), &claims, func(*jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, apperror.NewErrUnauthorized()
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, apperror.NewErrUnauthorized()
	}

	canConsume := false
	for _, permission := range permissions {
		if slices.Contains(claims.Permissions, permission) {
			canConsume = true
			break
		}
	}

	if !canConsume {
		return nil, apperror.NewErrUnauthorized()
	}

	return &AuthenticatorToken{
		UserID:      userID,
		Permissions: claims.Permissions,
	}, nil
}

func NewLocalAuthenticator(
	logger logger.Logger,
) Authenticator {
	return &localAuthenticator{
		logger: logger,
	}
}
