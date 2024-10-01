package usecases

import (
	"card-payment-api/internal/domain/dtos"
	"card-payment-api/internal/domain/entities"
	"card-payment-api/internal/ports"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

type Tokenizer interface {
	CreateToken(tokenRequest dtos.TokenDTO) (string, error)
	GetToken(token string) (*entities.TokenEntity, error)
}

// TokenUseCase is a use case for creating tokens
type TokenUseCase struct {
	tokenRepository ports.TokenRepository
}

// NewTokenUseCase creates a new TokenUseCase
func NewTokenUseCase(
	tokenRepository ports.TokenRepository,
) *TokenUseCase {
	return &TokenUseCase{
		tokenRepository: tokenRepository,
	}
}

// CreateToken creates a new token
func (t *TokenUseCase) CreateToken(tokenRequest dtos.TokenDTO) (string, error) {
	token, err := generateNumericToken(tokenRequest)
	if err != nil {
		return "", err
	}
	err = t.tokenRepository.Save(
		&entities.TokenEntity{
			Token:          token,
			CardNumber:     tokenRequest.CardNumber,
			CardHolderName: tokenRequest.CardHolderName,
			ExpirationDate: tokenRequest.ExpirationDate,
			CVV:            tokenRequest.CVV,
		},
	)
	if err != nil {
		return "", err

	}

	return token, nil
}

func (t *TokenUseCase) GetToken(token string) (*entities.TokenEntity, error) {
	tokenEntity, err := t.tokenRepository.Get(token)
	if err != nil {
		return nil, err
	}
	return tokenEntity, nil
}

// generate deterministic token
func generateNumericToken(tokenRequest dtos.TokenDTO) (string, error) {
	// concat card info
	cardInfo := tokenRequest.CardNumber + tokenRequest.CardHolderName +
		tokenRequest.ExpirationDate + tokenRequest.CVV

	// generate hash
	hash := sha256.Sum256([]byte(cardInfo))

	// convert hash to hex string
	hashString := hex.EncodeToString(hash[:])

	// take first 20 characters
	token := ""
	for i := 0; i < 20; i++ {
		num, _ := strconv.ParseInt(string(hashString[i]), 16, 64)
		token += strconv.Itoa(int(num % 10))
	}
	return token, nil
}
