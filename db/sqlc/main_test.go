package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *pgxpool.Pool
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

// Главный метод тестирования всего функционала (точка входа всех юнит тестов пакета db)
func TestMain(m *testing.M) {
	var err error
	testDB, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("SQL conn error:", err)
	}
	defer testDB.Close() // Close the connection pool when the test is done

	testQueries = New(testDB) //Add queries to DB connection

	os.Exit(m.Run())

}
