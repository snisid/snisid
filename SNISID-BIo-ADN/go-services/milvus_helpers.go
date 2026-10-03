package bio_adn

import (
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// varcharColumnData safely extracts string data from a Milvus search-result
// column. It follows the same contract as
// services/afis-svc/internal/repository/milvus/vector_repo.go: VarChar
// columns returned by the SDK are *entity.ColumnVarChar.
func varcharColumnData(col entity.Column) ([]string, error) {
	if col == nil {
		return nil, fmt.Errorf("column not present in search results")
	}
	vc, isVarChar := col.(*entity.ColumnVarChar)
	if !isVarChar {
		return nil, fmt.Errorf("unexpected column type %T, expected *entity.ColumnVarChar", col)
	}
	return vc.Data(), nil
}

// at returns the element at index i, or an empty string when out of range.
func at(data []string, i int) string {
	if i < 0 || i >= len(data) {
		return ""
	}
	return data[i]
}
