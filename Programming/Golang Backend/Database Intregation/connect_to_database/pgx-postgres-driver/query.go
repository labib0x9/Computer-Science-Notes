package main

/*
	pgx  -  POSTGRES DRIVER
*/

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// CREATE TABLE urls (
//     id  SERIAL PRIMARY KEY,
//     url TEXT         NOT NULL
// );

type Url struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

func main() {
	//// ---- Connection ---- ////
	// pgURL := "postgres://username:password@localhost:5432/database_name"
	pgURL := "postgres://short:secret@localhost:5432/short"
	conn, err := pgx.Connect(context.Background(), pgURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	///// ---- QUERY ONE ---- /////
	/* For single row */

	query := "select * from urls where id=$1"
	row := conn.QueryRow(context.Background(), query, 10000) // 10000 not in the table

	// Not found error
	var foundURL Url
	if err := row.Scan(&foundURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			println("No rows returned")
		} else {
			panic(err)
		}
	} else {
		fmt.Println("Found=", foundURL)
	}

	///// ---- QUERY TWO ---- //////
	/* For Multiple Rows */

	query = "select * from urls"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	foundURLs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Url])
	if err != nil {
		panic("Row collection: " + err.Error())
	}

	fmt.Println("Found Bro=", foundURLs)

	///// ----- QUERY THREE ------ //////

}
