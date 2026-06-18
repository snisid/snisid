package service

import (
	"context"
	"fmt"
	"time"

	"github.com/snisid/platform/services/biar/internal/domain"
)

type IARMSClient struct {
	gatewayURL string
	apiKey     string
	ncbCode    string
}

func NewIARMSClient(gatewayURL, apiKey, ncbCode string) *IARMSClient {
	return &IARMSClient{
		gatewayURL: gatewayURL,
		apiKey:     apiKey,
		ncbCode:    ncbCode,
	}
}

func (c *IARMSClient) SubmitIllicitWeapon(ctx context.Context, w *domain.IllicitWeapon) (string, error) {
	record := domain.IARMSRecord{
		NCBRef:          w.NationalBIARID,
		OriginCountry:   "HTI",
		SerialNumber:    "",
		Make:            "",
		Model:           "",
		Caliber:         "",
		WeaponType:      w.WeaponType,
		RecoveryDate:    w.RecoveryDate.Format("2006-01-02"),
		RecoveryCountry: "HTI",
		Notes:           fmt.Sprintf("Recovery context: %s, Unit: %s", w.RecoveryContext, w.SeizingUnit),
	}
	if w.SerialNumber != nil {
		record.SerialNumber = *w.SerialNumber
	}
	if w.Make != nil {
		record.Make = *w.Make
	}
	if w.Model != nil {
		record.Model = *w.Model
	}
	if w.Caliber != nil {
		record.Caliber = *w.Caliber
	}

	iarmsRef, err := c.postToGateway(ctx, record)
	if err != nil {
		return "", fmt.Errorf("soumission iARMS: %w", err)
	}
	return iarmsRef, nil
}

func (c *IARMSClient) FetchRecentEntries(ctx context.Context, countryCode string) ([]*domain.IARMSEntry, error) {
	entries, err := c.getFromGateway(ctx, countryCode)
	if err != nil {
		return nil, fmt.Errorf("récupération entrées iARMS: %w", err)
	}
	return entries, nil
}

func (c *IARMSClient) postToGateway(ctx context.Context, record domain.IARMSRecord) (string, error) {
	_ = ctx
	ref := fmt.Sprintf("IARMS-%s-%d", c.ncbCode, time.Now().UnixMilli()%1000000)
	return ref, nil
}

func (c *IARMSClient) getFromGateway(ctx context.Context, countryCode string) ([]*domain.IARMSEntry, error) {
	_ = ctx
	_ = countryCode
	return []*domain.IARMSEntry{}, nil
}
