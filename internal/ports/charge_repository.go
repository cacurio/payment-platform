package ports

import "card-payment-api/internal/domain/entities"

type ChargeRepository interface {
	Save(charge *entities.ChargeEntity) error
	Get(id string) (*entities.ChargeEntity, error)
}
