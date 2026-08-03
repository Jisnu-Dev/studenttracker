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
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Department   DepartmentType `json:"department"`
	Semester     int            `json:"semester"`
	Age          int            `json:"age"`
	CreatedAtUTC time.Time      `json:"createdAtUtc"`
	UpdatedAtUTC time.Time      `json:"updatedAtUtc"`
}

type PatchStudent struct {
	Name       *string         `json:"name"`
	Email      *string         `json:"email"`
	Department *DepartmentType `json:"department"`
	Semester   *int            `json:"semester"`
	Age        *int            `json:"age"`
}
