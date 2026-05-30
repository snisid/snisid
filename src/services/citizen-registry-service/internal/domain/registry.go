package domain

import (
	"errors"
)

var (
	ErrCitizenNotFound = errors.New("citizen not found")
)

// CitizenView represents the read-optimized denormalized view of a citizen.
type CitizenView struct {
	NIU          string `json:"niu"`
	FullName     string `json:"fullName"`
	DateOfBirth  string `json:"dateOfBirth"`
	PlaceOfBirth string `json:"placeOfBirth"`
	Status       string `json:"status"`
	LastUpdated  string `json:"lastUpdated"`
}
