package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"example.com/expenses-tracker/internal/models"
	"example.com/expenses-tracker/internal/repositories"
	"example.com/expenses-tracker/internal/requests"
	"golang.org/x/crypto/bcrypt"
)

const (
	errInvalidToken           = "invalid token"
	errFailedToCheckToken     = "failed to check token: %w"
	errFailedToGetUserByEmail = "failed to get user by email address: %w"
	errInvalidCredentials     = "user credentials incorrect"
)

type AuthHandler struct {
	userTokenRepository repositories.UserAuthRepository
	userRepository      repositories.UserRepository
}

// NewAuthHandler creates a new auth handler for checking an authenticated user
func NewAuthHandler(userTokenRepository repositories.UserAuthRepository, userRepository repositories.UserRepository) *AuthHandler {
	return &AuthHandler{
		userTokenRepository: userTokenRepository,
		userRepository:      userRepository,
	}
}

func (h *AuthHandler) HandleLoginRequest(ctx context.Context, request *requests.LoginRequest) (*models.User, error) {
	user, err := h.userRepository.GetUserByEmailAddress(ctx, request.EmailAddress)
	if err != nil {
		return nil, fmt.Errorf(errFailedToGetUserByEmail, err)
	}

	if user == nil {
		return nil, errors.New(errInvalidCredentials)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errors.New(errInvalidCredentials)
	}

	token, err := h.generateSecureTokenForUser(32)
	if err != nil {
		return user, err
	}

	h.persistTokenForUser(token, user)

	return user, nil
}

func (h *AuthHandler) persistTokenForUser(token string, user *models.User) (bool, error) {
	return true, nil
}

func (h *AuthHandler) generateSecureTokenForUser(tokenSize int64) (string, error) {
	bytes := make([]byte, tokenSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}

// ValidateToken will check that a specified token value is valid.
func (a *AuthHandler) ValidateToken(ctx context.Context, token string) (bool, error) {
	if token == "" || len(token) == 0 {
		return false, errors.New(errInvalidToken)
	}

	userToken, err := a.userTokenRepository.GetByAuthToken(ctx, token)
	if err != nil {
		return false, fmt.Errorf(errFailedToCheckToken, err)
	}

	if userToken == nil {
		return false, nil
	}

	return userToken.IsTokenValid(), nil
}
