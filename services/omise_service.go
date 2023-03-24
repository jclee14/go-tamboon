package services

import (
	"go-tamboon/config"
	"go-tamboon/models"
	"log"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type IOmiseService interface {
	CreateCharge(models.OmiseChargeParameter) (*models.OmiseChargeResponse, error)
}

type omiseService struct {
	omiseClient *omise.Client
}

func NewOmiseService(cfg config.IConfig) IOmiseService {
	client, err := omise.NewClient(cfg.OmisePublicKey(), cfg.OmiseSecretKey())
	if err != nil {
		log.Fatal(err)
	}
	return &omiseService{
		omiseClient: client,
	}
}

func (s *omiseService) CreateCharge(data models.OmiseChargeParameter) (*models.OmiseChargeResponse, error) {
	chargeData := &operations.CreateCharge{
		Card:     data.Card,
		Amount:   int64(data.Amount),
		Currency: data.Currency,
	}
	err := s.omiseClient.Do(&omise.Charge{}, chargeData)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
