package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Fungsi untuk membuat JWT
func GenerateToken(userID uint, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	// Buat claims 
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token berlaku 72 jam
		"iat":     time.Now().Unix(),
	}

	// Buat token baru dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

// Claims custom untuk menyimpan data user
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Fungsi untuk memvalidasi JWT
func ValidateToken(tokenString string) (*JWTClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Metode signing tidak terduga")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err 
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Token tidak valid")
}