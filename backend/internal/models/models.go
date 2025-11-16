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

type UserPreferences struct {
	Login string    `json:"login" db:"login"`
	Noon  time.Time `json:"noon" db:"noon"`
	Lang  string    `json:"lang" db:"lang"`
}

type UserCommonItem struct {
	UUID        uuid.UUID  `json:"uuid" db:"uuid"`
	Login       string     `json:"login" db:"login"`
	Path        string     `json:"path" db:"path"`
	Name        string     `json:"name" db:"name"`
	ProductUUID *uuid.UUID `json:"product_uuid,omitempty" db:"product_uuid"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

