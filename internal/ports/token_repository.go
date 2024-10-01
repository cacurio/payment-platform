package ports

import "card-payment-api/internal/domain/entities"

type TokenRepository interface {
	Save(token *entities.TokenEntity) error
	Get(token string) (*entities.TokenEntity, error)
}
