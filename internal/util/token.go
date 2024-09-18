package util

import (
	"context"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/domain"
	"github.com/Archetarcher/go-musthave-diploma-tpl.git/internal/logger"
	"github.com/go-chi/jwtauth/v5"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
	"time"
)

func CreateToken(user *domain.User, tokenConfig config.Token) (string, error) {
	logger.Log.Info("user id token", zap.Any("user", user))
	claims := jwt2.MapClaims{
		"id": user.ID,
	}

	jwtauth.SetExpiry(claims, time.Now().Add(time.Minute*time.Duration(tokenConfig.ExpiresInMinutes)))
	_, token, err := tokenConfig.AuthToken.Encode(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}
func IsAuthorized(requestContext context.Context) bool {
	token, _, err := jwtauth.FromContext(requestContext)

	if err != nil {
		logger.Log.Info("401 error token")
		return false
	}

	if token != nil && jwt.Validate(token) == nil {
		logger.Log.Info("token valid")
		return true
	}

	return false
}

func GetIdFromToken(requestContext context.Context) (int, error) {
	_, claims, err := jwtauth.FromContext(requestContext)
	if err != nil {
		return 0, err
	}
	userIdFromClaims := int(claims["id"].(float64))

	return userIdFromClaims, nil

}
func GenerateAuthToken(appConfig *config.AppConfig) {
	appConfig.Token.AuthToken = jwtauth.New("HS256", []byte(appConfig.Token.Key), nil)
}
