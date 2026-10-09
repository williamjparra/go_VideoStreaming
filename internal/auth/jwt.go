package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/joho/godotenv"
)

// Cargar variables de entorno y almacenarla en jwtSecret
var jwtSecret []byte

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No se encontraron variables de entorno: %v", err)
	}
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"username": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, "", errors.New("token inválido o expirado")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("no se pudieron obtener los claims del token")
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, "", errors.New("ID de usuario no válido dentro del token")
	}

	username, _ := claims["username"].(string)

	return int64(userIDFloat), username, nil
}