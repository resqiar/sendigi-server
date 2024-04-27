package utils

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(id string) string {
	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{"id": id},
	)

	JWT_SECRET := os.Getenv("JWT_SECRET")
	token, err := claims.SignedString([]byte(JWT_SECRET))
	if err != nil {
		log.Println("ERROR GENERATING JWT:", err)
		return ""
	}

	return token
}

func ParseJWT(token string) string {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Println("Failed to parse:", err)
		return ""
	}

	// Extract the claims from the token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Failed to extract claims:", err)
		return ""
	}

	userID := claims["id"].(string)
	if userID == "" {
		return ""
	}

	return userID
}
