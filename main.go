package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/api"
	db "github.com/techschool/simplebank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var (
	store         *db.Store
	conn          *pgxpool.Pool
	serverAddress = "localhost:8080"
)

func main() {
	var err error
	conn, err = pgxpool.New(context.Background(), dbSource) //Подключение к БД
	if err != nil {
		log.Fatal("SQL conn error:", err)
	}
	fmt.Println("*****CONNECTED****")
	defer conn.Close()
	store = db.NewStore(conn) //К запросам добавили подключение к БД
	server := api.NewServer(store)
	fmt.Println("Server started at 8080 port")
	if err := server.Start(serverAddress); err != nil {
		log.Fatal("Server cannot start", err)
	}

}
