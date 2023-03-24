package models

import "github.com/omise/omise-go"

type TransactionPayload struct {
	Headers []string
	Data    []string
	Row     int
}

type OmiseChargeParameter struct {
	Amount   int
	Currency string
	Card     string
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
