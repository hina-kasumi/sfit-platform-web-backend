package services

import (
	"fmt"
	"os"
	"sfit-platform-web-backend/entities"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshTokenService struct {
	expToken  int64
	secretKey string
}

func NewRefreshTokenService() *RefreshTokenService {
	tokenExpiration := os.Getenv("REFRESH_EXPARIATION")
	if tokenExpiration == "" {
		fmt.Println("REFRESH_EXPARIATION is not set")
		os.Exit(1)
	}
	exp, err := strconv.ParseInt(tokenExpiration, 10, 64)
	if err != nil {
		fmt.Println("JWT_EXPIRATION is not set or invalid")
		os.Exit(1)
	}

	secretKey := os.Getenv("REFRESH_SECRET")
	if secretKey == "" {
		fmt.Println("REFRESH_SECRET is not set")
		os.Exit(1)
	}

	return &RefreshTokenService{
		expToken:  exp,
		secretKey: secretKey,
	}
}

func (refreshSer *RefreshTokenService) GetRefreshTokenExp() int64 {
	return refreshSer.expToken
}

func (refreshSer *RefreshTokenService) GenerateRefreshToken(user entities.Users) (string, error) {
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
