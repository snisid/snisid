package indexes

import (
	"context"
	"fmt"

	"github.com/snisid/platform/services/bio-adn/pkg/models"
)

type DNAIndex struct {
	db models.Database
}

func NewDNAIndex(db models.Database) *DNAIndex {
	return &DNAIndex{db: db}
}

func (idx *DNAIndex) Create(ctx context.Context, profile *models.DNAProfile) error {
	if profile.SpecimenNumber == "" {
		return fmt.Errorf("specimen_number is required")
	}
	if profile.LociHash == "" {
		return fmt.Errorf("loci_hash is required")
	}
	return idx.db.CreateDNAProfile(ctx, profile)
}

func (idx *DNAIndex) GetByHash(ctx context.Context, hash string) (*models.DNAProfile, error) {
	return idx.db.GetDNAProfileByHash(ctx, hash)
}

func (idx *DNAIndex) GetBySpecimen(ctx context.Context, specimen string) (*models.DNAProfile, error) {
	return idx.db.GetDNAProfileBySpecimen(ctx, specimen)
}

func (idx *DNAIndex) SearchByIndexType(ctx context.Context, indexType string, limit, offset int) ([]models.DNAProfile, int, error) {
	return idx.db.SearchDNAProfiles(ctx, indexType, limit, offset)
}

func (idx *DNAIndex) GetUnuploaded(ctx context.Context, level string) ([]models.DNAProfile, error) {
	results, err := idx.db.GetUnuploadedDNAProfiles(ctx, level)
	if err != nil {
		return nil, err
	}
	var profiles []models.DNAProfile
	for _, r := range results {
		profiles = append(profiles, models.DNAProfile{
			SampleID:       r["id"].(string),
			SpecimenNumber: r["specimen_number"].(string),
			IndexType:      r["index_type"].(string),
		})
	}
	return profiles, nil
}

func (idx *DNAIndex) MarkUploaded(ctx context.Context, id, level string) error {
	return idx.db.MarkUploaded(ctx, id, level)
}

func (idx *DNAIndex) Expunge(ctx context.Context, id string) error {
	return idx.db.MarkExpunged(ctx, id)
}
