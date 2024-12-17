package service

import (
	"database/sql"
	"elkonterm/config"
	"elkonterm/internal/models"
	"elkonterm/internal/repository"
	"elkonterm/pkg/helper"
	"elkonterm/pkg/logger"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type authService struct {
	repo   *repository.Repository
	logger *logger.Logger
	cfg    *config.Config
}

type jwtCustomClaim struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
	RoleId uuid.UUID `json:"role_id"`
	Type   string    `json:"type"`
}

func newAuthService(repo *repository.Repository, logger *logger.Logger, cfg *config.Config) *authService {
	return &authService{
		repo:   repo,
		logger: logger,
		cfg:    cfg,
	}
}

func (s *authService) CreateToken(user models.User, tokenType string, expiresAt time.Time) (*models.Token, error) {
	claims := &jwtCustomClaim{
		UserId: user.ID,
		RoleId: user.Role.ID,
		Type:   tokenType,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(config.GetConfig().JWTSecret))
	if err != nil {
		return nil, serviceError(err, codes.Internal)
	}

	return &models.Token{
		Value:     token,
		Type:      tokenType,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *authService) GenerateTokens(user models.User) (*models.Token, *models.Token, error) {
	accessExpiresAt := time.Now().Add(time.Duration(s.cfg.JWTAccessExpirationHours) * time.Hour)
	refreshExpiresAt := time.Now().Add(time.Duration(s.cfg.JWTRefreshExpirationDays) * time.Hour * 24)

	accessToken, err := s.CreateToken(user, config.TokenTypeAccess, accessExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.CreateToken(user, config.TokenTypeRefresh, refreshExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) ParseToken(token string) (*jwtCustomClaim, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwtCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*jwtCustomClaim)
	if !ok {
		return nil, errors.New("token claims are not of type *jwtCustomClaim")
	}

	return claims, nil
}

func (s *authService) LoginAdmin(input models.LoginRequest) (*models.Token, *models.Token, error) {
	user, err := s.repo.User.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, serviceError(errWrongEmailOrPassword, codes.InvalidArgument)
		}
		return nil, nil, serviceError(err, codes.Internal)
	}

	roleID := user.RoleID.String()

	if roleID != config.SuperAdminRoleID && roleID != config.ManagerRoleID {
		return nil, nil, serviceError(errWrongEmailOrPassword, codes.InvalidArgument)
	}

	hashPassword, err := helper.GenerateHash(s.cfg, input.Password)
	if err != nil {
		return nil, nil, serviceError(err, codes.Internal)
	}

	if user.Password != hashPassword {
		return nil, nil, serviceError(errWrongEmailOrPassword, codes.Unauthenticated)
	}

	return s.GenerateTokens(models.User{
		ID: user.ID,
		Role: models.Role{
			ID: user.RoleID,
		},
	})
}
