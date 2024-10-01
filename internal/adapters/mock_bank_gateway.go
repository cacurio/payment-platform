package adapters

import (
	"bytes"
	"card-payment-api/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type MockBankGateway struct {
}

func NewMockBankGateway() *MockBankGateway {
	return &MockBankGateway{}
}

func (m *MockBankGateway) ProcessPayment(cardNumber, expiryMonth, expiryYear, currency string, amount float64) (string, error) {
	// implement api client
	bankUrl := os.Getenv("BANK_URL")
	apiUrl := bankUrl + "/process-payment"
	apiBody := map[string]interface{}{
		"cardNumber":  cardNumber,
		"expiryMonth": expiryMonth,
		"expiryYear":  expiryYear,
		"amount":      amount,
		"currency":    currency,
	}

	convertToJson, err := json.Marshal(apiBody)
	if err != nil {
		return "", err
	}
	fmt.Println(string(convertToJson))
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(convertToJson))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", domain.ErrNoProcessPayment
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", domain.ErrNoProcessPayment
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", domain.ErrNoProcessPayment
	}

	paymentNumber := response["authorizationCode"].(string)

	return paymentNumber, nil

}
func (m *MockBankGateway) RefundPayment(paymentId, currency string, amount float64) (string, error) {
	bankUrl := os.Getenv("BANK_URL")
	// implement api client
	apiUrl := bankUrl + "/process-refund"
	fmt.Println(apiUrl)
	apiBody := map[string]interface{}{
		"originalTransactionId": paymentId,
		"amount":                amount,
		"currency":              currency,
	}

	convertToJson, err := json.Marshal(apiBody)
	if err != nil {
		return "", err
	}
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(convertToJson))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("refund failed")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	refundId := response["refundId"].(string)

	return refundId, nil

}
