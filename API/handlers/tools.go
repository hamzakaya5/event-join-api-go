package handlers

import (
	"fmt"
	"gamegos_case/models"
	"math/rand"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Password does not match")
		return false
	}

	return err == nil
}

var jwtKey = []byte("my_secret_key")

// 1. Generate JWT Token
func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func Rand1to70() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(70-1+1) + 1 // => 1..70
}

func GetLock(key string) *sync.Mutex {
	models.LocksMu.Lock()
	defer models.LocksMu.Unlock()

	if l, ok := models.Locks[key]; ok {
		return l
	}
	l := &sync.Mutex{}
	models.Locks[key] = l
	return l
}
