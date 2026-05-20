package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mike-testut/task-api/internal/models"
	"github.com/mike-testut/task-api/internal/store"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret, ok = os.LookupEnv("JWT_SECRET")
var mySigningKey = []byte(jwtSecret)

type AuthService struct {
	userStore store.UserStore
}

func NewAuthService(us store.UserStore) *AuthService {
	return &AuthService{userStore: us}
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(username, password string) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return models.User{}, fmt.Errorf("could not hash password: %v", err)
	}

	user, err := s.userStore.CreateUser(username, string(hash))
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userStore.GetUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("login failed: invalid credentials: %v", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("login failed: invalid credentials: %v", err)
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %v", err)
	}
	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
