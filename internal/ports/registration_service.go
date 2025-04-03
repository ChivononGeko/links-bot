package ports

import "certificate/internal/domain"

type RegistrationService interface {
	GenerateUniqueLink(baseURL string) (string, error)
	RegisterUser(token, name, phone, birthday string) error
	GetTokenUsage(token string) (*domain.TokenUsage, error)
	GetUsedTokens() ([]domain.Registration, error)
	GetUnusedTokens() ([]domain.Registration, error)
	ValidateAndDecode(encryptedToken string) (string, error)
}
