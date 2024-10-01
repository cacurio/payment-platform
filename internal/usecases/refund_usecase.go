package usecases

import (
	"card-payment-api/internal/domain/dtos"
	"card-payment-api/internal/ports"
	"errors"
	"fmt"
)

type RefundUseCase struct {
	BankGateway      ports.BankGateway
	ChargeRepository ports.ChargeRepository
}

func NewRefundUseCase(bankGateway ports.BankGateway, chargeRepository ports.ChargeRepository) *RefundUseCase {
	return &RefundUseCase{
		BankGateway:      bankGateway,
		ChargeRepository: chargeRepository,
	}
}

func (r *RefundUseCase) Execute(refund dtos.RefundDTO) (string, error) {

	charge, err := r.ChargeRepository.Get(refund.PaymentId)
	if err != nil {
		return "", err
	}

	if charge.Status == "approved" {
		refundResponse, err := r.BankGateway.RefundPayment(refund.PaymentId, charge.Currency, refund.Amount)
		if err != nil {
			return "", err
		}
		charge.Status = "refunded"
		charge.RefundId = refundResponse
		err = r.ChargeRepository.Save(charge)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		return refundResponse, nil
	} else {
		return "", errors.New("refund not allowed")
	}

}
