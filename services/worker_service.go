package services

import (
	"errors"
	"fmt"
	"go-tamboon/config"
	"go-tamboon/models"
	"strconv"
	"strings"
)

type IWorkerService interface {
	ConsumePayload(ch <-chan models.TransactionPayload, doneCh <-chan struct{})
}

type workerService struct {
	config       config.IConfig
	omiseService IOmiseService

	totalReceived       float64
	successfullyDonated float64
	faultyDonation      float64
}

func NewWorkerService(config config.IConfig, omiseService IOmiseService) IWorkerService {
	return &workerService{
		config:              config,
		omiseService:        omiseService,
		totalReceived:       0,
		successfullyDonated: 0,
		faultyDonation:      0,
	}
}

func (s *workerService) readFile(filePath string) error {
	return nil
}

func (s *workerService) updateTotalReceived() {

}

func (s *workerService) updateFaultyDonation(errorAmount float64) {
	s.faultyDonation += errorAmount
}

func (s *workerService) ConsumePayload(ch <-chan models.TransactionPayload, doneCh <-chan struct{}) {
selectLoop:
	for {
		select {
		case payload := <-ch:
			fmt.Printf("%v\n", payload)
			s.doDonate(payload)
		case <-doneCh:
			break selectLoop
		}
	}
}

func (s *workerService) doDonate(payload models.TransactionPayload) error {
	if len(payload.Headers) == 0 || len(payload.Data) == 0 {
		return errors.New("invalid data")
	}

	inputData := map[string]string{}
	for idx, header := range payload.Headers {
		formatHeader := strings.ToLower(header)
		inputData[formatHeader] = ""
		if len(payload.Data) <= idx {
			continue
		}
		inputData[formatHeader] = payload.Data[idx]
	}

	createTokenParams, err := s.createTokenParameter(inputData)
	if err != nil {
		return err
	}

	createTokenResp, err := s.omiseService.CreateToken(createTokenParams)
	if err != nil {
		return err
	}

	createChargeParams, err := s.createChargeParameter(inputData, createTokenResp.ID)
	if err != nil {
		return err
	}

	_, err = s.omiseService.CreateCharge(createChargeParams)
	if err != nil {
		return err
	}

	// TODO: Collect result

	return nil
}

func (s *workerService) createTokenParameter(data map[string]string) (models.OmiseCreateTokenParameter, error) {
	res := models.OmiseCreateTokenParameter{}

	name, ok := data["name"]
	if !ok {
		return res, errors.New("name is invalid")
	}
	res.Name = name

	ccnumber, ok := data["ccnumber"]
	if !ok {
		return res, errors.New("ccnumber is invalid")
	}
	res.Number = ccnumber

	cvv, ok := data["cvv"]
	if !ok {
		return res, errors.New("cvv is invalid")
	}
	res.SecurityCode = cvv

	expMonth, err := strconv.Atoi(data["expmonth"])
	if err != nil {
		return res, err
	}
	res.ExpirationMonth = expMonth

	expYear, err := strconv.Atoi(data["expyear"])
	if err != nil {
		return res, err
	}
	res.ExpirationYear = expYear

	return res, nil
}

func (s *workerService) createChargeParameter(data map[string]string, token string) (models.OmiseChargeParameter, error) {
	res := models.OmiseChargeParameter{
		Currency: s.config.OmiseCurrency(),
	}

	if len(token) == 0 {
		return res, errors.New("token is invalid")
	}
	res.Card = token

	amount, err := strconv.Atoi(data["amountsubunits"])
	if err != nil {
		return res, err
	}
	res.Amount = amount

	return res, nil
}
