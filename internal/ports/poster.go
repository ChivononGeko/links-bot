package ports

import "certificate/internal/domain"

type PosterAPI interface {
	ChangeClientBonus(clientID int) error
	CreateClient(c domain.Client) (int, error)
}
