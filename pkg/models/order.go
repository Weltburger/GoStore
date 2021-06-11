package models

import (
	"GoStore/pkg/custom_types"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Order struct {
	UUID        uuid.UUID             `json:"uuid" db:"uuid"`
	ProductUUID uuid.UUID             `json:"product_uuid" db:"product_uuid"`
	Price       int                   `json:"price" db:"price"`
	Quantity    int                   `json:"quantity" db:"quantity"`
	Email       string                `json:"email" db:"email"`
	Status      int                   `json:"status" db:"status"`
	CreateAt    time.Time             `json:"create_at" db:"create_at"`
	UpdatedAt   time.Time             `json:"updated_at" db:"updated_at"`
	DeletedAt   custom_types.NullTime `json:"deleted_at" db:"deleted_at"`
}
