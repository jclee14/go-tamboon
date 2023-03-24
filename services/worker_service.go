package services

import (
	"errors"
	"fmt"
	"go-tamboon/config"
	"go-tamboon/models"
)

type IWorkerService interface {
	ConsumePayload(ch <-chan models.TransactionPayload, doneCh <-chan struct{})
}

type workerService struct {
	config config.IConfig

	totalReceived       float64
	successfullyDonated float64
	faultyDonation      float64
}

func NewWorkerService(config config.IConfig) IWorkerService {
	return &workerService{
		config:              config,
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

	data := map[string]string{}
	for idx, header := range payload.Headers {
		data[header] = ""
		if len(payload.Data) <= idx {
			continue
		}
		data[header] = payload.Data[idx]
	}

	return nil
}

func (s *workerService) createCharge(payload models.TransactionPayload) error {

	return nil
}
