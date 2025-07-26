package jwtservice

import (
	"fmt"
	"os"
	"sfit-platform-web-backend/entities"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetRefreshTokenExp() int64 {
	tokenExpiration := os.Getenv("REFRESH_EXPARIATION")
	if tokenExpiration == "" {
		fmt.Println("JWT_EXPIRATION is not set")
		os.Exit(1)
	}
	exp, err := strconv.ParseInt(tokenExpiration, 10, 64)
	if err != nil {
		fmt.Println("JWT_EXPIRATION is not set or invalid")
		os.Exit(1)
	}
	return exp
}

func getRefreshSecretKey() string {
	secretKey := os.Getenv("REFRESH_SECRET")
	if secretKey == "" {
		fmt.Println("JWT_SECRET is not set")
		os.Exit(1)
	}
	return secretKey
}

func GenerateRefreshToken(user entities.Users) (string, error) {
	secretKey := []byte(getRefreshSecretKey())
	expSecs := GetRefreshTokenExp()

	exp := time.Now().Unix() + expSecs

	claims := jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": user.ID.String(),
		"exp": exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseRefreshToken(token string) (jwt.Claims, error) {
	secretKey := []byte(getRefreshSecretKey())
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
