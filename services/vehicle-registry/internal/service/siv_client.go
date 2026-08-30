package service

import "fmt"

type SIVClient struct {
	baseURL string
}

func NewSIVClient() *SIVClient {
	return &SIVClient{baseURL: "https://siv.api.example.com"}
}

func (c *SIVClient) CheckTechnicalInspection(vin string) (bool, error) {
	if vin == "" {
		return false, fmt.Errorf("empty VIN")
	}
	return true, nil
}

func (c *SIVClient) GetInsuranceStatus(plate string) (string, error) {
	if plate == "" {
		return "", fmt.Errorf("empty plate")
	}
	return "valid", nil
}
