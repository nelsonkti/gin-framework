package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenData struct {
	UserId   string
	ExpireAt int64
}

// ParseToken 解析token
func ParseToken(token string, secretKey string) (*TokenData, error) {
	newToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := newToken.Claims.(jwt.MapClaims)
	if !ok || !newToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	UserId, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id type")
	}

	expireAt, ok := claims["expire_at"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid expire_at type")
	}

	return &TokenData{
		UserId:   UserId,
		ExpireAt: int64(expireAt),
	}, nil
}

// GenerateToken 生成token
func GenerateToken(userId string, expiration time.Duration, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userId,
		"expire_at": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
