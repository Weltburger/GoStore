package models

import (
	"GoStore/pkg/custom_types"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Product struct {
	UUID        uuid.UUID `json:"uuid" db:"uuid"`
	User        *User     `json:"user"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       int       `json:"price" db:"price"`
	Quantity    int       `json:"quantity" db:"quantity"`
	//Comments    []*Comment            `json:"comments"`
	CreateAt  time.Time             `json:"create_at" db:"create_at"`
	UpdatedAt time.Time             `json:"updated_at" db:"updated_at"`
	DeletedAt custom_types.NullTime `json:"deleted_at" db:"deleted_at"`
}
