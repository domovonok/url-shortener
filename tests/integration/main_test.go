package integration

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"

	"github.com/domovonok/url-shortener/tests/testdb"
)

var ConnectionString string

func TestMain(m *testing.M) {
	ctx := context.Background()

	err := godotenv.Load("../../.env.integration")
	if err != nil {
		log.Fatal(".env.integration file not found:", err)
	}

	pgContainer, dbConnectionString, err := testdb.TestWithMigrations()
	if err != nil {
		log.Fatal("container creation error:", err)
	}
	ConnectionString = dbConnectionString

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatal("failed to terminate container:", err)
		}
	}()

	exitCode := m.Run()
	os.Exit(exitCode)
}
