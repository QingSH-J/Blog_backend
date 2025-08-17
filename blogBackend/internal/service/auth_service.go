package service

import (
	"fmt"
	"log"
	"os"
	"project/internal/model"
	"project/internal/store"
	"time"

	"github.com/golang-jwt/jwt/v5" // Assuming you are using this package for JWT token generation
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(name string, password string, email string) (*model.UserLog, error)
	Login(email string, password string) (string, error)
}

type authService struct {
	userStore store.UserStore
}

func NewAuthService(userStore store.UserStore) AuthService {
	return &authService{userStore: userStore}
}

func (s *authService) Register(name string, password string, email string) (*model.UserLog, error) {
	//check the email is already registered
	existingUser, err := s.userStore.FindUserByEmail(email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}
	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	newUser := &model.UserLog{
		UserName: name,
		Email:    email,
		Password: string(hashedPassword),
	}
	if err := s.userStore.CreateUser(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	return newUser, nil
}

func (s *authService) Login(email string, password string) (string, error) {
	log.Printf("Login attempt for email: %s", email)

	//Find the exsiting user by email
	exsitingUser, err := s.userStore.FindUserByEmail(email)
	if err != nil {
		log.Printf("Failed to find user by email %s: %v", email, err)
		return "", fmt.Errorf("failed to find user: %v", err)
	}
	if exsitingUser == nil {
		log.Printf("User not found for email: %s", email)
		return "", fmt.Errorf("user not found")
	}

	log.Printf("User found: %s, checking password...", exsitingUser.Email)

	//Check the password
	if err := bcrypt.CompareHashAndPassword([]byte(exsitingUser.Password), []byte(password)); err != nil {
		log.Printf("Password verification failed for user %s: %v", email, err)
		return "", fmt.Errorf("invalid password")
	}

	log.Printf("Login successful for user: %s", email)
	// Generate a token using the JWT_SECRET from environment
	claims := jwt.MapClaims{
		"user_id": exsitingUser.ID,
		"email":   exsitingUser.Email,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用环境变量中的JWT密钥，与中间件保持一致
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set in environment variables")
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}
