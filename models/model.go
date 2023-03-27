package models

import "github.com/omise/omise-go"

type TransactionPayload struct {
	Headers []string
	Data    []string
	Row     int
}

type InputData struct {
}

type OmiseCreateTokenParameter struct {
	Name            string
	Number          string
	ExpirationMonth int
	ExpirationYear  int
	SecurityCode    string
}

type OmiseChargeParameter struct {
	Amount   int
	Currency string
	Card     string
}

type OmiseCreateTokenResponse struct {
	omise.Card
}

type OmiseChargeResponse struct {
	omise.Charge

	// Object         string    `json:"object"`
	// Id             string    `json:"id"`
	// Amount         int       `json:"amount"`
	// Net            int       `json:"net"`
	// Fee            int       `json:"fee"`
	// FeeVat         int       `json:"fee_vat"`
	// Currency       string    `json:"currency"`
	// FailureCode    string    `json:"failure_code"`
	// FailureMessage string    `json:"failure_message"`
	// Status         string    `json:"status"`
	// CreatedAt      time.Time `json:"created_at"`
	// PaidAt         time.Time `json:"paid_at"`
	// ExpiresAt      time.Time `json:"expires_at"`
}

type OmiseError struct {
	Location   string `json:"location"`
	StatusCode int    `json:"status"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (err *OmiseError) Error() string

func (err *OmiseError) String() string
