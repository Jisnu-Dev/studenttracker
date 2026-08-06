package models

import (
	"time"
)

type Admin struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name" validate:"required,min=2,max=100,name"`
	Email        string    `json:"email" validate:"required,email,max=254"`
	Password     string    `json:"password" validate:"required,min=8,max=64,password"`
	CreatedAtUTC time.Time `json:"createdAtUtc"`
	UpdatedAtUTC time.Time `json:"updatedAtUtc"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
