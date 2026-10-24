package main

/*
	pgx  -  POSTGRES DRIVER, You must ask yourself why ???? prevent using pooling.

	labib0x9@TestGoCode $ go run main.go
	GOROUTINE= 0 ERRor= conn busy
	GOROUTINE= 1 ERRor= conn busy
	GOROUTINE= 2 ERRor= conn busy
	GOROUTINE= 3 FOUND ID= 3 URL= https://labib0x0hunter.netlify.app/#projects
*/

import (
	"context"
	"fmt"
	"sync"

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
	conn, err := pgx.Connect(context.Background(), pgURL) // 1-TCP connection.
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	var wg sync.WaitGroup
	for i := range 4 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			query := "select * from urls where id=$1"
			row := conn.QueryRow(context.Background(), query, id)
			var found Url
			if err := row.Scan(&found.Id, &found.Url); err != nil {
				fmt.Println("GOROUTINE=", id, "ERRor=", err)
			} else {
				fmt.Println("GOROUTINE=", id, "FOUND ID=", found.Id, "URL=", found.Url)
			}
		}(i)
	}

	wg.Wait()
}
