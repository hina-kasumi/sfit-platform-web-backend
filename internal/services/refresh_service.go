package services

import (
	"fmt"
	"sfit-platform-web-backend/internal/config"
	"sfit-platform-web-backend/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshTokenService struct {
	expToken  int64
	secretKey string
}

func NewRefreshTokenService(cfg config.RefreshTokenConfig) *RefreshTokenService {
	return &RefreshTokenService{
		expToken:  int64(cfg.ExpirationSec),
		secretKey: cfg.SecretKey,
	}
}

func (refreshSer *RefreshTokenService) GetRefreshTokenExp() int64 {
	return refreshSer.expToken
}

func (refreshSer *RefreshTokenService) GenerateRefreshToken(user model.Users) (string, error) {
	secretKey := []byte(refreshSer.secretKey)
	expSecs := refreshSer.expToken

	exp := time.Now().Unix() + expSecs

	claims := jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": user.ID.String(),
		"exp": exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (refreshSer *RefreshTokenService) ParseRefreshToken(token string) (jwt.Claims, error) {
	secretKey := []byte(refreshSer.secretKey)
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken.Claims, nil
}
