package tokenService

import (
	"github.com/hryt430/task-management/internal/modules/auth/domain/entity"
	"github.com/hryt430/task-management/pkg/token"
)

type TokenUseCase interface {
	GenerateAccessToken(user *entity.User) (string, error)
	GenerateRefreshToken(user *entity.User) (string, error)
	ValidateAccessToken(tokenString string) (*token.Claims, error)
	RevokeAccessToken(tokenString string) error
	IsTokenRevoked(tokenString string) bool
}
