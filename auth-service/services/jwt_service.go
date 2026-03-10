package services

import (
	"strconv"
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	Secret                   string
	AccessExpirationMinutes  int
	RefreshExpirationHours   int
}

type TokenClaims struct {
	CredentialID int64  `json:"credential_id"`
	EmployeeID   int64  `json:"employee_id"`
	Email        string `json:"email"`
	TokenType    string `json:"token_type"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, accessMinutes string, refreshHours string) *JWTService {
	accessExp, err := strconv.Atoi(accessMinutes)
	if err != nil {
		accessExp = 15
	}

	refreshExp, err := strconv.Atoi(refreshHours)
	if err != nil {
		refreshExp = 168
	}

	return &JWTService{
		Secret:                  secret,
		AccessExpirationMinutes: accessExp,
		RefreshExpirationHours:  refreshExp,
	}
}

func (s *JWTService) GenerateAccessToken(credentialID int64, employeeID int64, email string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(s.AccessExpirationMinutes) * time.Minute)

	claims := TokenClaims{
		CredentialID: credentialID,
		EmployeeID:   employeeID,
		Email:        email,
		TokenType:    "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(employeeID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Secret))
}

func (s *JWTService) GenerateRefreshToken(credentialID int64, employeeID int64, email string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(s.RefreshExpirationHours) * time.Hour)

	claims := TokenClaims{
		CredentialID: credentialID,
		EmployeeID:   employeeID,
		Email:        email,
		TokenType:    "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(employeeID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Secret))
}

func (s *JWTService) ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, errors.New("provided token is not an access token")
	}

	return claims, nil
}