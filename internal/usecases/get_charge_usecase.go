package usecases

import (
	"card-payment-api/internal/domain/entities"
	"card-payment-api/internal/ports"
)

type GetChargeUseCase struct {
	chargeRepository ports.ChargeRepository
}

func NewGetChargeUseCase(chargeRepository ports.ChargeRepository) *GetChargeUseCase {
	return &GetChargeUseCase{
		chargeRepository: chargeRepository,
	}
}

func (g *GetChargeUseCase) GetCharge(id string) (*entities.ChargeEntity, error) {
	return g.chargeRepository.Get(id)
}
