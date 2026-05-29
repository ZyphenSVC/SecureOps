package auth

import (
	"context"
	"errors"
	"time"

	"github.com/ZyphenSVC/SecureOps/backend/internal/audit"
	"github.com/ZyphenSVC/SecureOps/backend/internal/rbac"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/ZyphenSVC/SecureOps/backend/internal/users"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	users      *users.Repository
	rbac       *rbac.Repository
	audit      *audit.Repository
	jwtSecret  []byte
	bcryptCost int
}

func NewService(
	users *users.Repository,
	rbacRepo *rbac.Repository,
	auditRepo *audit.Repository,
	jwtSecret string,
	bcryptCost int,
) *Service {
	return &Service{
		users:      users,
		rbac:       rbacRepo,
		audit:      auditRepo,
		jwtSecret:  []byte(jwtSecret),
		bcryptCost: bcryptCost,
	}
}

func (s *Service) Register(ctx context.Context, email, password, fullName string) (*users.User, string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return nil, "", err
	}

	user, err := s.users.Create(ctx, email, string(passwordHash), fullName)
	if err != nil {
		return nil, "", err
	}

	resourceID := user.ID
	_ = s.audit.Record(ctx, &user.ID, "user.registered", "user", &resourceID, "{}")

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*users.User, string, error) {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	resourceID := user.ID
	_ = s.audit.Record(ctx, &user.ID, "auth.login", "user", &resourceID, "{}")

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *Service) GenerateToken(user *users.User) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   now.Unix(),
		"exp":   now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.jwtSecret)
}
