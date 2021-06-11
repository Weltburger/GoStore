package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	UUID      uuid.UUID             `json:"uuid" db:"uuid"`
	Name      string                `json:"name" db:"name"`
	Email     string                `json:"email" db:"email"`
	Password  string                `json:"-" db:"password"`
	CreateAt  time.Time             `json:"create_at" db:"create_at"`
	UpdatedAt time.Time             `json:"updated_at" db:"updated_at"`
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func CreateToken(length int) (string, error) {
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}

func Auth(JWTKey []byte, value string) (string, error) {
	token, err := jwt.Parse(value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return JWTKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenClaim := claims["token"].(string)
		return tokenClaim, nil
	}

	return "", errors.New("bad access token")
}

