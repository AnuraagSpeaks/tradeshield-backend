package token

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("payshield-secret-jwt-key-2026-production")

func getSigningKey() []byte {
	if envSecret := os.Getenv("JWT_SECRET"); envSecret != "" {
		return []byte(envSecret)
	}
	return SecretKey
}

type Claims struct {
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, email, role, orgID string) (string, error) {
	claims := Claims{
		UserID:         userID,
		Email:          email,
		Role:           role,
		OrganizationID: orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "payshield-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSigningKey())
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getSigningKey(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
