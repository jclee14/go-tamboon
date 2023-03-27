package services

import (
	"go-tamboon/config"
	"go-tamboon/models"
	"time"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type IOmiseService interface {
	CreateToken(data models.OmiseCreateTokenParameter) (*models.OmiseCreateTokenResponse, error)
	CreateCharge(models.OmiseChargeParameter) (*models.OmiseChargeResponse, error)
}

type omiseService struct {
	omiseClient *omise.Client
}

func NewOmiseService(cfg config.IConfig) (IOmiseService, error) {
	client, err := omise.NewClient(cfg.OmisePublicKey(), cfg.OmiseSecretKey())
	if err != nil {
		return nil, err
	}
	return &omiseService{
		omiseClient: client,
	}, nil
}

func (s *omiseService) CreateToken(data models.OmiseCreateTokenParameter) (*models.OmiseCreateTokenResponse, error) {
	createToken := &operations.CreateToken{
		Name:            data.Name,
		Number:          data.Number,
		ExpirationMonth: time.Month(data.ExpirationMonth),
		ExpirationYear:  data.ExpirationYear,
		SecurityCode:    data.SecurityCode,
	}
	result := &models.OmiseCreateTokenResponse{}

	err := s.omiseClient.Do(result, createToken)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *omiseService) CreateCharge(data models.OmiseChargeParameter) (*models.OmiseChargeResponse, error) {
	chargeData := &operations.CreateCharge{
		Card:     data.Card,
		Amount:   int64(data.Amount),
		Currency: data.Currency,
	}
	chargeResult := &models.OmiseChargeResponse{}

	if err := s.omiseClient.Do(chargeResult, chargeData); err != nil {
		return nil, err
		// if omiseErr, ok := err.(*models.OmiseError); ok {
		// 	log.Fatal(omiseErr.Code + " " + omiseErr.Message)
		// } else {
		// 	return nil, err
		// }
	}

	return chargeResult, nil
}
