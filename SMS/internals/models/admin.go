package models

import (
	"time"
)

type Admin struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name" binding:"required"`
	Email        string    `json:"email" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	CreatedAtUTC time.Time `json:"createdAtUtc"`
	UpdatedAtUTC time.Time `json:"updatedAtUtc"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
