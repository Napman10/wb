package domain

import (
	"time"

	"github.com/google/uuid"
)

type Gender int

const (
	Male Gender = iota
	Female
)

type Employee struct {
	ID       uuid.UUID `json:"id"`
	Fullname string    `json:"fullname"`
	Gender   Gender    `json:"gender"`
	Age      uint      `json:"age"`
	Email    string    `json:"email"`
	Address  string    `json:"address"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
