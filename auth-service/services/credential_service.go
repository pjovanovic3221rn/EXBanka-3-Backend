package services

import (
	"errors"
	"strings"
	"time"

	"auth-service/models"
	"auth-service/repository"
)

type CredentialService struct {
	CredentialRepo *repository.CredentialRepository
	JWTService     *JWTService
}

func NewCredentialService(credentialRepo *repository.CredentialRepository, jwtService *JWTService) *CredentialService {
	return &CredentialService{
		CredentialRepo: credentialRepo,
		JWTService:     jwtService,
	}
}

func (s *CredentialService) CreateCredential(req *models.CreateCredentialRequest) (*models.Credential, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.EmployeeID <= 0 {
		return nil, errors.New("employee_id must be greater than 0")
	}

	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	if !strings.Contains(req.Email, "@") {
		return nil, errors.New("invalid email format")
	}

	exists, err := s.CredentialRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	existingByEmployee, err := s.CredentialRepo.GetByEmployeeID(req.EmployeeID)
	if err == nil && existingByEmployee != nil {
		return nil, errors.New("credential for this employee already exists")
	}

	activationToken, err := GenerateRandomToken(32)
	if err != nil {
		return nil, err
	}

	credential := &models.Credential{
		EmployeeID:      req.EmployeeID,
		Email:           req.Email,
		IsActive:        req.IsActive,
		ActivationToken: &activationToken,
	}

	err = s.CredentialRepo.Create(credential)
	if err != nil {
		return nil, err
	}

	return credential, nil
}

func (s *CredentialService) ActivateAccount(token string, password string, confirmPassword string) error {

	if token == "" {
		return errors.New("activation token is required")
	}

	if password != confirmPassword {
		return errors.New("passwords do not match")
	}

	err := ValidatePassword(password)
	if err != nil {
		return err
	}

	credential, err := s.CredentialRepo.GetByActivationToken(token)
	if err != nil {
		return errors.New("invalid activation token")
	}

	if credential.IsActive {
		return errors.New("account already activated")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	salt, err := GenerateSalt(16)
	if err != nil {
		return err
	}

	err = s.CredentialRepo.ActivateAccount(credential.ID, hash, salt)
	if err != nil {
		return err
	}

	return nil
}

func (s *CredentialService) Login(email string, password string) (*models.LoginResponse, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	credential, err := s.CredentialRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !credential.IsActive {
		return nil, errors.New("account is not activated")
	}

	if credential.PasswordHash == "" {
		return nil, errors.New("account has no password set")
	}

	if !CheckPasswordHash(password, credential.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.JWTService.GenerateAccessToken(
		credential.ID,
		credential.EmployeeID,
		credential.Email,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(
		credential.ID,
		credential.EmployeeID,
		credential.Email,
	)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

func (s *CredentialService) RefreshAccessToken(refreshToken string) (*models.RefreshResponse, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	claims, err := s.JWTService.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("provided token is not a refresh token")
	}

	credential, err := s.CredentialRepo.GetByEmail(claims.Email)
	if err != nil {
		return nil, errors.New("credential not found")
	}

	if !credential.IsActive {
		return nil, errors.New("account is not active")
	}

	newAccessToken, err := s.JWTService.GenerateAccessToken(
		credential.ID,
		credential.EmployeeID,
		credential.Email,
	)
	if err != nil {
		return nil, err
	}

	return &models.RefreshResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
	}, nil
}

func (s *CredentialService) ForgotPassword(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return "", errors.New("email is required")
	}

	credential, err := s.CredentialRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("credential not found")
	}

	if !credential.IsActive {
		return "", errors.New("account is not active")
	}

	resetToken, err := GenerateRandomToken(32)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(1 * time.Hour)

	err = s.CredentialRepo.SetResetToken(credential.ID, resetToken, expiresAt)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *CredentialService) ResetPassword(resetToken string, password string, confirmPassword string) error {
	if resetToken == "" {
		return errors.New("reset token is required")
	}

	if password != confirmPassword {
		return errors.New("passwords do not match")
	}

	err := ValidatePassword(password)
	if err != nil {
		return err
	}

	credential, err := s.CredentialRepo.GetByResetToken(resetToken)
	if err != nil {
		return errors.New("invalid reset token")
	}

	if credential.ResetTokenExpires == nil || time.Now().After(*credential.ResetTokenExpires) {
		return errors.New("reset token has expired")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	salt, err := GenerateSalt(16)
	if err != nil {
		return err
	}

	err = s.CredentialRepo.ResetPassword(credential.ID, hash, salt)
	if err != nil {
		return err
	}

	return nil
}