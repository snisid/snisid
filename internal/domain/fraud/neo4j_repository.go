package fraud

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphRepository interface {
	AddIdentityNode(ctx context.Context, identityID, agency string) error
	CheckFraudRing(ctx context.Context, identityID string) (bool, error)
}

type neo4jRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jRepository(driver neo4j.DriverWithContext) GraphRepository {
	return &neo4jRepository{driver: driver}
}

func (r *neo4jRepository) AddIdentityNode(ctx context.Context, identityID, agency string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MERGE (i:Identity {id: $id}) SET i.agency = $agency`
		params := map[string]any{
			"id":     identityID,
			"agency": agency,
		}
		return tx.Run(ctx, query, params)
	})
	return err
}

func (r *neo4jRepository) CheckFraudRing(ctx context.Context, identityID string) (bool, error) {
	// A simple mock query for fraud ring detection
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (i:Identity {id: $id})-[:SHARES_DEVICE]->(d:Device)<-[:SHARES_DEVICE]-(other:Identity)
			RETURN count(other) > 2 AS isFraudRing
		`
		params := map[string]any{"id": identityID}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return false, err
		}
		if result.Next(ctx) {
			return result.Record().Values[0].(bool), nil
		}
		return false, nil
	})
	if err != nil {
		return false, err
	}
	return res.(bool), nil
}
