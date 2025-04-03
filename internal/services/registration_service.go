package services

import (
	"certificate/internal/domain"
	"certificate/internal/ports"
	"fmt"
	"time"
)

type RegistrationService struct {
	repo      ports.RegistrationRepository
	posterAPI ports.PosterAPI
}

func NewRegistrationService(repo ports.RegistrationRepository, posterAPI ports.PosterAPI) *RegistrationService {
	return &RegistrationService{
		repo:      repo,
		posterAPI: posterAPI,
	}
}

// Генерация уникальной ссылки на основе времени
func (s *RegistrationService) GenerateUniqueLink(baseURL string) (string, error) {
	token := fmt.Sprintf("%d", time.Now().UnixNano())

	err := s.repo.Create(token)
	if err != nil {
		return "", err
	}

	encryptedToken, err := encryptToken(token)
	if err != nil {
		return "", err
	}

	return baseURL + encryptedToken, nil
}

// Проверка, был ли уже использован токен
func (s *RegistrationService) ValidateAndDecode(encryptedToken string) (string, error) {
	decodedToken, err := decryptToken(encryptedToken)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// Проверяем, что токен существует и не был использован
	reg, err := s.repo.GetByToken(decodedToken)
	if err != nil || reg.Used {
		return "", fmt.Errorf("invalid or used token")
	}

	return decodedToken, nil
}

// Пометить токен как использованный
func (s *RegistrationService) MarkTokenUsed(token string, client domain.Client) error {
	return s.repo.MarkTokenUsed(token, client.Name, client.Phone)
}

// Выполнить полную регистрацию клиента
func (s *RegistrationService) RegisterUser(token, name, phone, birthday string) error {
	client := domain.Client{
		Name:     name,
		Phone:    phone,
		Birthday: birthday,
	}

	clientID, err := s.posterAPI.CreateClient(client)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Начисляем бонусы
	if err := s.posterAPI.ChangeClientBonus(clientID); err != nil {
		return fmt.Errorf("failed to change client bonus: %w", err)
	}

	// Помечаем токен как использованный
	if err := s.MarkTokenUsed(token, client); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}

// Получить информацию пользователя, который использовал токен
func (s *RegistrationService) GetTokenUsage(token string) (*domain.TokenUsage, error) {
	usage, err := s.repo.GetTokenUsage(token)
	if err != nil {
		return nil, err
	}

	return usage, nil
}

// Получить список использованных токенов
func (s *RegistrationService) GetUsedTokens() ([]domain.Registration, error) {
	return s.repo.GetUsedTokens()
}

// Получить список не использованных токенов
func (s *RegistrationService) GetUnusedTokens() ([]domain.Registration, error) {
	return s.repo.GetUnusedTokens()
}
