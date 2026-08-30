package service

import "fmt"

type FoVesClient struct {
	baseURL string
}

func NewFoVesClient() *FoVesClient {
	return &FoVesClient{baseURL: "https://foves.api.example.com"}
}

func (c *FoVesClient) VerifyRegistration(vin string) (bool, error) {
	if vin == "" {
		return false, fmt.Errorf("empty VIN")
	}
	return true, nil
}

func (c *FoVesClient) ReportTheft(plate string) error {
	if plate == "" {
		return fmt.Errorf("empty plate")
	}
	return nil
}
