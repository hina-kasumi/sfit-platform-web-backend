package jwtservice

import (
	"fmt"
	"os"
	"sfit-platform-web-backend/entities"
	redisservice "sfit-platform-web-backend/services/redis_service"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func getBlackListPrefix() string {
	blacklistPrefix := os.Getenv("JWT_BLACKLIST_PREFIX")
	if blacklistPrefix == "" {
		fmt.Println("JWT_BLACKLIST_PREFIX is not set")
		os.Exit(1)
	}
	return blacklistPrefix
}

func GetTokenExpiration() int64 {
	tokenExpiration := os.Getenv("JWT_EXPIRATION")
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

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		fmt.Println("JWT_SECRET is not set")
		os.Exit(1)
	}
	return secretKey
}

func GenerateToken(user entities.Users) (string, error) {
	secretKey := []byte(getSecretKey())
	expSecs := GetTokenExpiration()

	exp := time.Now().Unix() + expSecs

	claims := jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": user.ID.String(),
		"exp": exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(token string) error {
	_, err := ParseToken(token)
	if err != nil {
		return err
	}
	inBlacklist, err := IsTokenBlacklisted(token)
	if err != nil {
		return fmt.Errorf("error checking token blacklist: %v", err)
	}
	if inBlacklist {
		return fmt.Errorf("token is blacklisted")
	}
	return nil
}
func ParseToken(token string) (jwt.Claims, error) {
	secretKey := []byte(getSecretKey())
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

func IsTokenBlacklisted(token string) (bool, error) {
	blacklistKey, err := GenBlacklistKey(token)
	if err != nil {
		return false, err
	}
	redisClient, err := redisservice.GetRedisValue(blacklistKey)
	if err != nil {
		if err == redis.Nil {
			return false, nil // Token is not blacklisted
		}
		return false, err // Error occurred while checking blacklist
	}
	if redisClient != nil {
		return true, nil // Token is blacklisted
	}
	return false, nil // Token is not blacklisted
}

func BlacklistToken(token string) error {
	blacklistKey, err := GenBlacklistKey(token)
	if err != nil {
		return err
	}
	err = redisservice.SetRedisExpire(blacklistKey, "", GetTokenExpiration())
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %v", err)
	}
	return nil
}

func GenBlacklistKey(token string) (string, error) {
	blacklistPrefix := getBlackListPrefix()
	claims, err := ParseToken(token)
	if err != nil {
		return "", err
	}
	sub, err := claims.GetSubject()
	if err != nil {
		return "", err
	}
	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("unable to assert claims as MapClaims")
	}
	jit, ok := mapClaims["jti"]
	if !ok {
		return "", fmt.Errorf("jti claim not found in token")
	}
	blacklistKey := blacklistPrefix + sub + "_" + jit.(string)

	return blacklistKey, nil
}
