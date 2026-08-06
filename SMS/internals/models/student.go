package models

import (
	"time"
)

type DepartmentType string

const (
	CSE   DepartmentType = "CSE"
	IT    DepartmentType = "IT"
	ECE   DepartmentType = "ECE"
	EEE   DepartmentType = "EEE"
	MECH  DepartmentType = "MECH"
	CIVIL DepartmentType = "CIVIL"
	AIDS  DepartmentType = "AIDS"
	AIML  DepartmentType = "AIML"
)

type Student struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name" validate:"required,min=2,max=100,name"`
	Email        string         `json:"email" validate:"required,email,max=254"`
	Department   DepartmentType `json:"department" validate:"required,department"`
	Semester     int            `json:"semester" validate:"required,gte=1,lte=8"`
	Age          int            `json:"age" validate:"required,gte=18,lte=60"`
	CreatedAtUTC time.Time      `json:"createdAtUtc"`
	UpdatedAtUTC time.Time      `json:"updatedAtUtc"`
}

type PatchStudent struct {
	Name       *string         `json:"name" validate:"omitempty,min=2,max=100,name"`
	Email      *string         `json:"email" validate:"omitempty,email,max=254"`
	Department *DepartmentType `json:"department" validate:"omitempty,department"`
	Semester   *int            `json:"semester" validate:"omitempty,gte=1,lte=8"`
	Age        *int            `json:"age" validate:"omitempty,gte=18,lte=60"`
}
