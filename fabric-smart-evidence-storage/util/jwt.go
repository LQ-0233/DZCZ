package util

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// 生成 JWT Token
func GenerateJWT(username string, role string) (string, error) {

	// 生成 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// 使用密钥签名 JWT Token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	// 返回 JWT Token
	return tokenString, nil
}

func ParseJWT(tokenString string) (username string, role string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return "", "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid token")
	}
	username = claims["username"].(string)
	role = claims["role"].(string)
	return username, role, nil
}
