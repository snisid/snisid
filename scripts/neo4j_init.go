package main

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/snisid/platform/internal/platform/database"
	"github.com/snisid/platform/internal/platform/logger"
)

func main() {
	// ... (env loading)
	uri := "bolt://localhost:7687"
	user := "neo4j"
	pass := "neo4jpass"

	driver, err := database.NewNeo4jDriver(uri, user, pass)
	if err != nil {
		logger.Fatal("failed to connect to neo4j", err)
	}
	defer driver.Close(context.Background())

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	// Create Indexes for performance
	queries := []string{
		"CREATE INDEX identity_id IF NOT EXISTS FOR (i:Identity) ON (i.id)",
		"CREATE INDEX agency_name IF NOT EXISTS FOR (a:Agency) ON (a.name)",
		"CREATE CONSTRAINT identity_id_unique IF NOT EXISTS FOR (i:Identity) REQUIRE i.id IS UNIQUE",
	}

	for _, q := range queries {
		_, err := session.Run(context.Background(), q, nil)
		if err != nil {
			logger.Error("failed to create index", err)
		}
	}
	logger.Info("Neo4j indexes initialized")
}
