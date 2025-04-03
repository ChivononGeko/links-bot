package ports

import "certificate/internal/domain"

type RegistrationRepository interface {
	Create(token string) error
	GetByToken(token string) (*domain.Registration, error)
	MarkTokenUsed(token, name, phone string) error
	GetTokenUsage(token string) (*domain.TokenUsage, error)
	GetUsedTokens() ([]domain.Registration, error)
	GetUnusedTokens() ([]domain.Registration, error)
}
