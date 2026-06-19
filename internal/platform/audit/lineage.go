package audit

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/snisid/platform/internal/platform/database"
	"github.com/snisid/platform/internal/platform/logger"
)

func TrackLineage(ctx context.Context, driver neo4j.DriverWithContext, identityID, source, processor string) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
		MERGE (i:Identity {id: $id})
		MERGE (s:DataSource {name: $source})
		MERGE (p:Processor {name: $processor})
		MERGE (s)-[:PRODUCED]->(i)
		MERGE (i)-[:PROCESSED_BY]->(p)
		SET i.lineage_verified = true, p.timestamp = datetime()
	`
	_, err := session.Run(ctx, query, map[string]interface{}{
		"id":        identityID,
		"source":    source,
		"processor": processor,
	})
	if err != nil {
		logger.Error(ctx, "failed to track data lineage", err)
	}
}
