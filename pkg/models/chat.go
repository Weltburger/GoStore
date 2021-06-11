package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Chat struct {
	UUID      uuid.UUID `json:"uuid" db:"uuid"`
	Data      []Message `json:"data" db:"data"`
	CreateAt  time.Time `json:"create_at" db:"create_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
}
