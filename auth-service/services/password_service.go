package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func GenerateSalt(length int) (string, error) {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidatePassword(password string) error {

	if len(password) < 8 || len(password) > 32 {
		return errors.New("password must be between 8 and 32 characters")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`)
	hasLower := regexp.MustCompile(`[a-z]`)
	hasDigits := regexp.MustCompile(`[0-9].*[0-9]`)

	if !hasUpper.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !hasLower.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !hasDigits.MatchString(password) {
		return errors.New("password must contain at least two numbers")
	}

	return nil
}