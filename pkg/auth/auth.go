package auth

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

func GenerateAccessToken(accessTokenTTL time.Duration, userId int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   fmt.Sprintf("%d", userId),
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Role: role,
	})
	accessToken, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET_KEY")))
	if err != nil {
		return "", errors.New("can't generate access token")
	}
	return accessToken, nil
}

func ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("TOKEN_SECRET_KEY")), nil
	})
	if err != nil {
		return 0, "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, "", errors.New("token claims are not of type *CustomClaims")
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, "", err
	}
	return userId, claims.Role, nil
}

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("PASSWORD_HASH_SALT"))))
}
