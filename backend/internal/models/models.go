package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductDetails struct {
	UUID     uuid.UUID `json:"uuid" db:"uuid"`
	Name     string    `json:"name" db:"name"`
	Ccal     int64     `json:"ccal" db:"ccal"`
	Fats     int64     `json:"fats" db:"fats"`
	Proteins int64     `json:"proteins" db:"proteins"`
	Carbs    int64     `json:"carbs" db:"carbs"`
}

type Record struct {
	UUID        uuid.UUID `json:"uuid" db:"uuid"`
	ProductUUID uuid.UUID `json:"product_uuid" db:"product_uuid"`
	Amount      int64     `json:"amount" db:"amount"`
	Login       string    `json:"login" db:"login"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

