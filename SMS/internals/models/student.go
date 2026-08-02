package models

import (
	"time"
)

type departmentType string

const (
	CSE   departmentType = "CSE"
	IT    departmentType = "IT"
	ECE   departmentType = "ECE"
	EEE   departmentType = "EEE"
	MECH  departmentType = "MECH"
	CIVIL departmentType = "CIVIL"
	AIDS  departmentType = "AIDS"
	AIML  departmentType = "AIML"
)

type Student struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Department   departmentType `json:"department"`
	Semester     int            `json:"semester"`
	Age          int            `json:"age"`
	CreatedAtUTC time.Time      `json:"createdAtUtc"`
	UpdatedAtUTC time.Time      `json:"updatedAtUtc"`
}

type PatchStudent struct {
	Name       *string         `json:"name"`
	Email      *string         `json:"email"`
	Department *departmentType `json:"department"`
	Semester   *int            `json:"semester"`
	Age        *int            `json:"age"`
}
