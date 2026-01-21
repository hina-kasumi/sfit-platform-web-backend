package services

import (
	"context"
	"fmt"
	"sfit-platform-web-backend/internal/config"
	"sfit-platform-web-backend/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type JwtService struct {
	redisClient     *redis.Client
	ctx             context.Context
	blacklistPrefix string
	expToken        int64
	secretKey       string
}

func NewJwtService(cfg *config.JwtConfig, redisClient *redis.Client, ctx context.Context) *JwtService {
	return &JwtService{
		redisClient:     redisClient,
		ctx:             ctx,
		blacklistPrefix: "BLACKLIST_PREFIX_",
		expToken:        int64(cfg.ExpirationSec),
		secretKey:       cfg.SecretKey,
	}
}

func (jwtSer *JwtService) GenerateToken(user model.Users) (string, error) {
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
	_, err = jwtSer.getRedisValue(blacklistKey)
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
	err = jwtSer.setRedisExpire(blacklistKey, "", jwtSer.expToken)
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

func (jwtSer *JwtService) setRedisExpire(key string, value string, expiration int64) error {
	return jwtSer.redisClient.Set(jwtSer.ctx, key, value, time.Duration(expiration)*time.Second).Err()
}

func (jwtSer *JwtService) getRedisValue(key string) (string, error) {
	return jwtSer.redisClient.Get(jwtSer.ctx, key).Result()
}
