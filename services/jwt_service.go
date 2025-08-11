package services

import (
	"fmt"
	"os"
	"sfit-platform-web-backend/entities"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type JwtService struct {
	redisService    *RedisService
	blacklistPrefix string
	expToken        int64
	secretKey       string
}

func NewJwtService(redisService *RedisService) *JwtService {
	blacklistPrefix := os.Getenv("JWT_BLACKLIST_PREFIX")
	if blacklistPrefix == "" {
		fmt.Println("JWT_BLACKLIST_PREFIX is not set")
		os.Exit(1)
	}

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

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		fmt.Println("JWT_SECRET is not set")
		os.Exit(1)
	}

	return &JwtService{
		redisService:    redisService,
		blacklistPrefix: blacklistPrefix,
		expToken:        exp,
		secretKey:       secretKey,
	}
}

func (jwtSer *JwtService) GenerateToken(user entities.Users) (string, error) {
	secretKey := []byte(jwtSer.secretKey)
	expSecs := jwtSer.expToken
	exp := time.Now().Unix() + expSecs

	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = string(role.RoleID)
	}

	claims := jwt.MapClaims{
		"jti":   uuid.New().String(),
		"sub":   user.ID.String(),
		"exp":   exp,
		"roles": roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (jwtSer *JwtService) ParseToken(tokenStr string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSer.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid or non-MapClaims token")
	}

	sub, _ := claims.GetSubject()
	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, fmt.Errorf("jti claim not found or not a string")
	}

	blacklistKey := jwtSer.GenBlacklistKey(jti, sub)
	_, err = jwtSer.redisService.GetRedisValue(blacklistKey)
	if err == nil || err != redis.Nil {
		return nil, fmt.Errorf("token is blacklisted")
	}

	return claims, nil
}

func (jwtSer *JwtService) BlacklistToken(tokenStr string) error {
	claims, err := jwtSer.ParseToken(tokenStr)

	sub, _ := claims.GetSubject()
	jti, ok := claims["jti"].(string)
	if !ok {
		return fmt.Errorf("jti claim not found or not a string")
	}
	blacklistKey := jwtSer.GenBlacklistKey(jti, sub)
	if err != nil {
		return err
	}
	err = jwtSer.redisService.SetRedisExpire(blacklistKey, "", jwtSer.expToken)
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %v", err)
	}
	return nil
}

func (jwtSer *JwtService) GenBlacklistKey(jit string, sub string) string {
	blacklistPrefix := jwtSer.blacklistPrefix
	blacklistKey := blacklistPrefix + sub + "_" + jit

	return blacklistKey
}
