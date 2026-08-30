package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCitizenAlreadyExists = errors.New("citizen already exists")
	ErrInvalidData          = errors.New("invalid citizen data")
)

type BiometricStatus string

const (
	StatusPendingBiometrics BiometricStatus = "PENDING_BIOMETRICS"
	StatusActive            BiometricStatus = "ACTIVE"
	StatusSuspended         BiometricStatus = "SUSPENDED"
)

type Citizen struct {
	ID           uuid.UUID
	NIU          string // National Identification Number
	FirstName    string
	LastName     string
	DateOfBirth  string
	PlaceOfBirth string
	Status       BiometricStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewCitizen(firstName, lastName, dob, pob string) (*Citizen, error) {
	if firstName == "" || lastName == "" {
		return nil, ErrInvalidData
	}

	return &Citizen{
		ID:           uuid.New(),
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirth:  dob,
		PlaceOfBirth: pob,
		Status:       StatusPendingBiometrics,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}, nil
}
