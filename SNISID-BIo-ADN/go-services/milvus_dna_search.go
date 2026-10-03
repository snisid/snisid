package bio_adn

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	DNACollectionName = "snisid_dna_profiles"
	DNAEmbeddingDim   = 512
)

type DNAMilvusSearcher struct {
	client client.Client
}

func NewDNAMilvusSearcher(c client.Client) *DNAMilvusSearcher {
	return &DNAMilvusSearcher{client: c}
}

func (s *DNAMilvusSearcher) EnsureCollection(ctx context.Context) error {
	exists, err := s.client.HasCollection(ctx, DNACollectionName)
	if err != nil {
		return fmt.Errorf("check dna collection: %w", err)
	}
	if exists {
		return nil
	}

	schema := &entity.Schema{
		CollectionName: DNACollectionName,
		Description:    "SNISID DNA Profile Vector Store — CODIS-compatible embeddings",
		AutoID:         false,
		Fields: []*entity.Field{
			{
				Name:       "profile_id",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: true,
				TypeParams: map[string]string{"max_length": "128"},
			},
			{
				Name:       "niu",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "10"},
			},
			{
				Name:       "profile_type",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "50"},
			},
			{
				Name:       "case_reference",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "100"},
			},
			{
				Name:       "status",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "30"},
			},
			{
				Name:       "submitting_agency",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "100"},
			},
			{
				Name:     "dna_embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": fmt.Sprintf("%d", DNAEmbeddingDim),
				},
			},
			{
				Name:     "locus_scores",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "24",
				},
			},
		},
	}

	if err := s.client.CreateCollection(ctx, schema, entity.DefaultShardNumber); err != nil {
		return fmt.Errorf("create dna collection: %w", err)
	}

	idx, err := entity.NewIndexIvfFlat(entity.L2, 128)
	if err != nil {
		return fmt.Errorf("create ivf_flat index spec for dna_embedding: %w", err)
	}
	if err := s.client.CreateIndex(ctx, DNACollectionName, "dna_embedding", idx, false); err != nil {
		return fmt.Errorf("create dna embedding index: %w", err)
	}

	locusIdx, err := entity.NewIndexIvfFlat(entity.L2, 64)
	if err != nil {
		return fmt.Errorf("create ivf_flat index spec for locus_scores: %w", err)
	}
	if err := s.client.CreateIndex(ctx, DNACollectionName, "locus_scores", locusIdx, false); err != nil {
		return fmt.Errorf("create locus scores index: %w", err)
	}

	return nil
}

type DNAProfileVector struct {
	ProfileID        string
	NIU              string
	ProfileType      string
	CaseReference    string
	Status           string
	SubmittingAgency string
	DNAEmbedding     []float32
	LocusScores      []float32
}

type DNAMatchResult struct {
	ProfileID     string
	NIU           string
	CaseReference string
	Score         float32
	MatchType     string
}

func (s *DNAMilvusSearcher) UpsertProfile(ctx context.Context, profile *DNAProfileVector) error {
	profileIDCol := entity.NewColumnVarChar("profile_id", []string{profile.ProfileID})
	niuCol := entity.NewColumnVarChar("niu", []string{profile.NIU})
	typeCol := entity.NewColumnVarChar("profile_type", []string{profile.ProfileType})
	caseCol := entity.NewColumnVarChar("case_reference", []string{profile.CaseReference})
	statusCol := entity.NewColumnVarChar("status", []string{profile.Status})
	agencyCol := entity.NewColumnVarChar("submitting_agency", []string{profile.SubmittingAgency})
	dnaCol := entity.NewColumnFloatVector("dna_embedding", DNAEmbeddingDim, [][]float32{profile.DNAEmbedding})
	locusCol := entity.NewColumnFloatVector("locus_scores", 24, [][]float32{profile.LocusScores})

	_, err := s.client.Insert(ctx, DNACollectionName, "",
		profileIDCol, niuCol, typeCol, caseCol, statusCol, agencyCol, dnaCol, locusCol,
	)
	return err
}

func (s *DNAMilvusSearcher) SearchByEmbedding(ctx context.Context, queryVector []float32, topK int) ([]DNAMatchResult, error) {
	if topK <= 0 {
		topK = 10
	}

	searchParam, err := entity.NewIndexIvfFlatSearchParam(16)
	if err != nil {
		return nil, fmt.Errorf("create ivf_flat search param: %w", err)
	}

	results, err := s.client.Search(ctx, DNACollectionName,
		[]string{},
		"",
		[]string{"profile_id", "niu", "case_reference", "status"},
		[]entity.Vector{entity.FloatVector(queryVector)},
		"dna_embedding",
		entity.L2,
		topK,
		searchParam,
	)
	if err != nil {
		return nil, fmt.Errorf("search dna profiles: %w", err)
	}

	var matches []DNAMatchResult
	for _, result := range results {
		profileIDs, err := varcharColumnData(result.Fields.GetColumn("profile_id"))
		if err != nil {
			return nil, fmt.Errorf("read profile_id column: %w", err)
		}
		nius, err := varcharColumnData(result.Fields.GetColumn("niu"))
		if err != nil {
			return nil, fmt.Errorf("read niu column: %w", err)
		}
		cases, err := varcharColumnData(result.Fields.GetColumn("case_reference"))
		if err != nil {
			return nil, fmt.Errorf("read case_reference column: %w", err)
		}

		for i := 0; i < result.ResultCount; i++ {
			matches = append(matches, DNAMatchResult{
				ProfileID:     at(profileIDs, i),
				NIU:           at(nius, i),
				CaseReference: at(cases, i),
				Score:         result.Scores[i],
				MatchType:     classifyDNAMatch(result.Scores[i]),
			})
		}
	}

	return matches, nil
}

// SearchByLocusScores compares a 24-locus STR score vector against the
// locus_scores field. Same Milvus SDK v2 contract as services/afis-svc
// (entity.ColumnVarChar results, error-checked constructors).
func (s *DNAMilvusSearcher) SearchByLocusScores(ctx context.Context, locusScores []float32, topK int) ([]DNAMatchResult, error) {
	if len(locusScores) != 24 {
		return nil, fmt.Errorf("locus scores must contain exactly 24 values, got %d", len(locusScores))
	}
	if topK <= 0 {
		topK = 5
	}

	searchParam, err := entity.NewIndexIvfFlatSearchParam(16)
	if err != nil {
		return nil, fmt.Errorf("create ivf_flat search param: %w", err)
	}

	results, err := s.client.Search(ctx, DNACollectionName,
		[]string{},
		"",
		[]string{"profile_id", "niu", "case_reference"},
		[]entity.Vector{entity.FloatVector(locusScores)},
		"locus_scores",
		entity.L2,
		topK,
		searchParam,
	)
	if err != nil {
		return nil, fmt.Errorf("search by locus: %w", err)
	}

	var matches []DNAMatchResult
	for _, result := range results {
		profileIDs, err := varcharColumnData(result.Fields.GetColumn("profile_id"))
		if err != nil {
			return nil, fmt.Errorf("read profile_id column: %w", err)
		}
		nius, err := varcharColumnData(result.Fields.GetColumn("niu"))
		if err != nil {
			return nil, fmt.Errorf("read niu column: %w", err)
		}
		cases, err := varcharColumnData(result.Fields.GetColumn("case_reference"))
		if err != nil {
			return nil, fmt.Errorf("read case_reference column: %w", err)
		}

		for i := 0; i < result.ResultCount; i++ {
			matches = append(matches, DNAMatchResult{
				ProfileID:     at(profileIDs, i),
				NIU:           at(nius, i),
				CaseReference: at(cases, i),
				Score:         result.Scores[i],
				MatchType:     classifyDNAMatch(result.Scores[i]),
			})
		}
	}

	return matches, nil
}

func classifyDNAMatch(score float32) string {
	switch {
	case score < 0.1:
		return "EXACT_MATCH"
	case score < 0.3:
		return "HIGH_CONFIDENCE"
	case score < 0.5:
		return "PARTIAL_MATCH"
	default:
		return "LOW_CONFIDENCE"
	}
}
