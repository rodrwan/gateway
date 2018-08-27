package gateway

import "time"

type Card struct {
	ID              string `json:"id,omitempty" db:"id"`
	UserID          string `json:"user_id,omitempty" db:"user_id"`
	ProductID       string `json:"product_id,omitempty" db:"product_id"`
	CardNumber      string `json:"card_number,omitempty" db:"pan"`
	ReferenceID     string `json:"reference_id,omitempty" db:"ref_id"`
	ReferenceEmail  string `json:"reference_email,omitempty" db:"ref_email"`
	ReferenceUserID string `json:"reference_user_id,omitempty" db:"ref_user_id"`
}

type CardQueryOptions struct {
	UserID string `q:"user_id"`
}

type CardQueryOption func(*CardQueryOptions)

func SetCardQueryOptions(src *CardQueryOptions) CardQueryOption {
	return func(dst *CardQueryOptions) {
		*dst = *src
	}
}

type CardService interface {
	Select(...CardQueryOption) ([]*Card, error)
	CardDeposits(string) ([]*CardDeposit, error)
}

type CardDeposit struct {
	ID        string    `json:"id,omitempty" db:"id"`
	Amount    int64     `json:"amount,omitempty" db:"amount"`
	PaymentID string    `json:"payment_id,omitempty" db:"payment_id"`
	Status    string    `json:"status,omitempty" db:"status"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	Fee       int64     `json:"fee,omitempty" db:"fee"`
	Total     int64     `json:"total,omitempty" db:"total"`
	Dollar    int64     `json:"dollar,omitempty" db:"usd"`
}
