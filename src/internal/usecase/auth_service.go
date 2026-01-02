package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pbmap_api/src/config"
	"pbmap_api/src/domain"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/repository"
	"pbmap_api/src/pkg/auth"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type AuthService interface {
	LoginWithSocial(ctx context.Context, req *dto.SocialLoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userUsecase    UserUsecase
	tokenRepo      repository.TokenRepository
	sessionRepo    repository.SessionRepository
	jwtService     *auth.JWTService
	googleClientID string
	lineChannelID  string
}

func NewAuthService(userUsecase UserUsecase, tokenRepo repository.TokenRepository, sessionRepo repository.SessionRepository, jwtService *auth.JWTService, cfg *config.Config) AuthService {
	return &authService{
		userUsecase:    userUsecase,
		tokenRepo:      tokenRepo,
		sessionRepo:    sessionRepo,
		jwtService:     jwtService,
		googleClientID: cfg.GoogleClientID,
		lineChannelID:  cfg.LineChannelID,
	}
}

func (s *authService) LoginWithSocial(ctx context.Context, req *dto.SocialLoginRequest) (*dto.LoginResponse, error) {
	var providerID string
	var err error

	switch req.Provider {
	case "google":
		providerID, err = s.verifyGoogleToken(req.AccessToken)
	case "line":
		providerID, err = s.verifyLineToken(req.AccessToken)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("provider verification failed: %v", err)
	}

	user, err := s.userUsecase.SyncUserFromSocial(ctx, req.Provider, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to sync user: %v", err)
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate app token: %v", err)
	}

	appTokenTTL := 72 * time.Hour
	if err := s.tokenRepo.SetAppToken(ctx, user.ID.String(), token, appTokenTTL); err != nil {
		return nil, fmt.Errorf("failed to store app token: %v", err)
	}

	// Generate and Store Refresh Token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	session := &domain.UserSession{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return &dto.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(appTokenTTL.Seconds()),
	}, nil
}

func (s *authService) verifyGoogleToken(token string) (string, error) {
	payload, err := idtoken.Validate(context.Background(), token, s.googleClientID)
	if err != nil {
		return "", fmt.Errorf("invalid google token: %v", err)
	}

	return payload.Subject, nil
}

func (s *authService) verifyLineToken(token string) (string, error) {
	apiURL := "https://api.line.me/oauth2/v2.1/verify"
	data := url.Values{}
	data.Set("id_token", token)
	data.Set("client_id", s.lineChannelID)

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return "", fmt.Errorf("failed to verify line id token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("line id token verification failed: status=%d, body=%s", resp.StatusCode, string(bodyBytes))
	}

	var result dto.LineIDTokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode line response: %v", err)
	}

	if result.Sub == "" {
		return "", fmt.Errorf("line id token valid but sub is empty")
	}

	return result.Sub, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.tokenRepo.DeleteAppToken(ctx, userID.String())
}

func (s *authService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	session, err := s.sessionRepo.GetSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = s.sessionRepo.RevokeSession(ctx, session.ID)
		return nil, fmt.Errorf("refresh token expired")
	}

	user, err := s.userUsecase.GetUser(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := s.sessionRepo.RevokeSession(ctx, session.ID); err != nil {
		return nil, fmt.Errorf("failed to revoke old session: %v", err)
	}

	newAccessToken, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	appTokenTTL := 72 * time.Hour
	if err := s.tokenRepo.SetAppToken(ctx, user.ID.String(), newAccessToken, appTokenTTL); err != nil {
		return nil, fmt.Errorf("failed to store app token: %v", err)
	}

	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	newSession := &domain.UserSession{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	if err := s.sessionRepo.CreateSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed to create new session: %v", err)
	}

	return &dto.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(appTokenTTL.Seconds()),
	}, nil
}

func (s *authService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
