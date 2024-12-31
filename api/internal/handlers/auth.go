package handlers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"example.com/expenses-tracker/api/internal/auth"
	"example.com/expenses-tracker/api/internal/repositories"
	"example.com/expenses-tracker/pkg/models"
	"example.com/expenses-tracker/pkg/requests"
	"example.com/expenses-tracker/pkg/responses"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("expired token")
	ErrInvalidCredentials = errors.New("user credentials incorrect")
)

const (
	errFailedToCheckToken      = "failed to check token: %w"
	errFailedToCreateToken     = "failed to create token: %w"
	errFailedToGetTokenByValue = "failed to get token by value: %w"
	errFailedToGetUserByEmail  = "failed to get user by email address: %w"
	errFailedToGetUserByToken  = "failed to get user by token: %w"
	errFailedToCreateSession   = "failed to create session: %w"
)

type AuthHandler struct {
	userSessionRepository repositories.UserSessionRepository
	userRepository        repositories.UserRepository
	tokenHandler          *auth.TokenHandler
}

// NewAuthHandler creates a new auth handler for checking an authenticated user
func NewAuthHandler(userTokenRepository repositories.UserSessionRepository, userRepository repositories.UserRepository, tokenHandler *auth.TokenHandler) *AuthHandler {
	return &AuthHandler{
		userSessionRepository: userTokenRepository,
		userRepository:        userRepository,
		tokenHandler:          tokenHandler,
	}
}

func (h *AuthHandler) HandleLoginRequest(ctx context.Context, digitalFingerprint string, request *requests.LoginRequest) (*responses.AuthenticatedUserResponse, error) {
	user, err := h.userRepository.GetUserByEmailAddress(ctx, request.EmailAddress)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	expiryDuration, _ := time.ParseDuration("20m")
	expiryTime := time.Now().Add(expiryDuration)
	sessionID, _ := bcrypt.GenerateFromPassword([]byte(user.ID), bcrypt.DefaultCost)
	fingerprintAsBytes, _ := bcrypt.GenerateFromPassword([]byte(digitalFingerprint), bcrypt.DefaultCost)

	session := &models.UserSession{
		ID:                 base64.RawStdEncoding.EncodeToString(sessionID),
		DigitalFingerPrint: base64.RawStdEncoding.EncodeToString(fingerprintAsBytes),
		SessionID:          user.ID,
		CreatedAt:          time.Now(),
		ExpiryTime:         expiryTime,
		User:               user,
	}

	token := h.tokenHandler.GenerateForSession(session, expiryTime)
	err = h.userSessionRepository.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf(errFailedToCreateSession, err)
	}

	response := &responses.AuthenticatedUserResponse{
		Token:      token,
		ExpiryTime: expiryTime,
	}

	return response, nil
}

func (h *AuthHandler) GetSessionIdFromToken(ctx context.Context, token string) (*string, error) {
	isValid, sessionId := h.tokenHandler.ValidateToken(token)
	if !isValid {
		return nil, ErrInvalidToken
	}

	session, err := h.userSessionRepository.GetBySessionID(ctx, *sessionId)
	if err != nil {
		return nil, err
	}

	if session.HasExpired() {
		return nil, ErrExpiredToken
	}

	return sessionId, nil
}

func (h *AuthHandler) GetUserBySessionID(ctx context.Context, token string) (*models.User, error) {
	return h.userRepository.GetUserBySessionID(ctx, token)
}

func (h *AuthHandler) HandleLogout(ctx context.Context, token string) (bool, error) {
	// TODO - blacklist a token

	return true, nil
}
