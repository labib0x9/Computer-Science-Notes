package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TableOP(dbConn *sqlx.DB, path string) {
	query, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	res, err := dbConn.Exec(string(query))
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func main() {
	user := "tempuser1"
	pass := "secret"
	host := "localhost"
	port := 5432
	dbName := "tempdb1"
	sslmode := "disable"

	dbSource := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s", user, pass, host, port, dbName, sslmode)

	fmt.Println(dbSource)

	dbConn, err := sqlx.Connect("postgres", dbSource)
	if err != nil {
		panic(err)
	}

	defer dbConn.Close()

	TableOP(dbConn, "./migrations/00001_create_users_table.down.sql")
	TableOP(dbConn, "./migrations/00002_create_email_verification_table.down.sql")
	TableOP(dbConn, "./migrations/00003_create_labs_table.down.sql")
}
