package models

import (
	"time"
)

type Admin struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	CreatedAtUTC time.Time `json:"createdAtUtc"`
	UpdatedAtUTC time.Time `json:"updatedAtUtc"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
