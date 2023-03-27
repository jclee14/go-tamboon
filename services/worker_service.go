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
	topDonors           []*struct {
		names  []string
		amount float64
	}
	sumDonors map[string]float64
}

func NewWorkerService(config config.IConfig, omiseService IOmiseService) IWorkerService {
	return &workerService{
		config:              config,
		omiseService:        omiseService,
		totalReceived:       0,
		successfullyDonated: 0,
		faultyDonation:      0,
		topDonors: make([]*struct {
			names  []string
			amount float64
		}, 3),
		sumDonors: make(map[string]float64, 0),
	}
}

func (s *workerService) updateTotalReceived(donateAmount float64) {
	s.totalReceived += donateAmount
}

func (s *workerService) updateSuccessDonation(donateAmount float64) {
	s.successfullyDonated += donateAmount
}

func (s *workerService) updateFaultyDonation(donateAmount float64) {
	s.faultyDonation += donateAmount
}

func (s *workerService) updateTopDonor(donorName string, donateAmount float64) {
	if _, ok := s.sumDonors[donorName]; ok {
		s.sumDonors[donorName] += donateAmount
	} else {
		s.sumDonors[donorName] = donateAmount
	}

	insertIdx := -1
	for idx, item := range s.topDonors {
		if s.sumDonors[donorName] > item.amount {
			insertIdx = idx
			break
		} else if s.sumDonors[donorName] == item.amount {
			item.names = append(item.names, donorName)
			break
		}
	}

	if insertIdx >= 0 {
		s.topDonors = append(s.topDonors[:insertIdx+1], s.topDonors[insertIdx:len(s.topDonors)]...)
		s.topDonors[insertIdx] = &struct {
			names  []string
			amount float64
		}{
			names:  []string{donorName},
			amount: s.sumDonors[donorName],
		}
	}

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

	// Collect amount
	donateAmount := getDonateAmountInFloat(inputData)
	if donateAmount < 1 {
		return errors.New("invalid donate amount")
	}
	s.updateTotalReceived(donateAmount)

	createTokenParams, err := s.createTokenParameter(inputData)
	if err != nil {
		s.updateFaultyDonation(donateAmount)
		return err
	}

	createTokenResp, err := s.omiseService.CreateToken(createTokenParams)
	if err != nil {
		s.updateFaultyDonation(donateAmount)
		return err
	}

	createChargeParams, err := s.createChargeParameter(inputData, createTokenResp.ID)
	if err != nil {
		s.updateFaultyDonation(donateAmount)
		return err
	}

	_, err = s.omiseService.CreateCharge(createChargeParams)
	if err != nil {
		s.updateFaultyDonation(donateAmount)
		return err
	}

	// TODO: Collect result
	s.updateSuccessDonation(donateAmount)
	// top donor

	return nil
}

func getDonateAmountInFloat(data map[string]string) float64 {
	amount, _ := strconv.ParseFloat(data[string(models.InputHeaderAmountsubunits)], 64)
	return amount
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

	amount, err := strconv.Atoi(data[string(models.InputHeaderAmountsubunits)])
	if err != nil {
		return res, err
	}
	res.Amount = amount

	return res, nil
}
