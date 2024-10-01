package usecases

import (
	"card-payment-api/internal/domain"
	"card-payment-api/internal/domain/dtos"
	"card-payment-api/internal/domain/entities"
	"card-payment-api/internal/ports"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Charger interface {
	CreateCharge(chargeRequest dtos.ChargeDTO) (*entities.ChargeEntity, error)
	GetCharge(id string) (*entities.ChargeEntity, error)
}

type ChargeUseCase struct {
	chargeRepository ports.ChargeRepository
	Tokenizer        Tokenizer
	BankGateway      ports.BankGateway
}

func NewChargeUseCase(
	chargeRepository ports.ChargeRepository,
	Tokenizer Tokenizer,
	BankGateway ports.BankGateway,
) *ChargeUseCase {
	return &ChargeUseCase{
		Tokenizer:        Tokenizer,
		BankGateway:      BankGateway,
		chargeRepository: chargeRepository,
	}
}

func (c *ChargeUseCase) CreateCharge(chargeRequest dtos.ChargeDTO) (*entities.ChargeEntity, error) {
	// get token
	tokenEntity, err := c.Tokenizer.GetToken(chargeRequest.Token)
	if errors.Is(err, domain.ErrInvalidToken) {
		return nil, err
	}

	// process payment with bank gateway
	processPaymentResponse, err := c.BankGateway.
		ProcessPayment(tokenEntity.CardNumber, tokenEntity.GetExpirationMonth(), tokenEntity.GetExpirationYear(), chargeRequest.Currency, chargeRequest.Amount)

	var chargeStatus string
	if errors.Is(err, domain.ErrNoProcessPayment) {
		chargeStatus = "declined"
	} else {
		chargeStatus = "approved"
	}

	// create id
	id := uuid.New().String()
	// save charge
	createdDate := time.Now().Format("2006-01-02 15:04:05")
	chargeEntity := &entities.ChargeEntity{
		ID:            id,
		Token:         tokenEntity.Token,
		Amount:        chargeRequest.Amount,
		Currency:      chargeRequest.Currency,
		Status:        chargeStatus,
		Created:       createdDate,
		PaymentNumber: processPaymentResponse,
	}

	err = c.chargeRepository.Save(chargeEntity)
	if errors.Is(err, domain.ErrNoCreatedCharge) {
		return nil, domain.ErrNoCreatedCharge
	}

	return chargeEntity, nil
}

func (c *ChargeUseCase) GetCharge(id string) (*entities.ChargeEntity, error) {
	return c.chargeRepository.Get(id)
}
