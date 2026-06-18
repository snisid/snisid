package service

import (
	"context"
	"math/rand"
)

type QualityService struct {
	minScore int16
}

func NewQualityService(minScore int16) *QualityService {
	return &QualityService{minScore: minScore}
}

func (s *QualityService) CheckQuality(ctx context.Context, imageBase64 string) (int16, error) {
	if len(imageBase64) == 0 {
		return 0, ErrEmptyImage
	}
	score := int16(50 + rand.Intn(51))
	return score, nil
}

func (s *QualityService) IsAcceptable(score int16) bool {
	return score >= s.minScore
}

func (s *QualityService) MinScore() int16 { return s.minScore }
