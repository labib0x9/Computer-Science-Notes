package main

/*
	pgx  -  POSTGRES DRIVER

	go get github.com/jackc/pgx/v5/pgxpool
*/

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

	// pool, err := pgxpool.New(context.Background(), pgURL)

	config, err := pgxpool.ParseConfig(pgURL)
	if err != nil {
		panic(err)
	}

	config.MinConns = 2
	config.MaxConns = 10
	config.MinIdleConns = 2
	config.MaxConnIdleTime = 2 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	var wg sync.WaitGroup
	for i := range 4 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			query := "select * from urls where id=$1"
			row := pool.QueryRow(context.Background(), query, id)
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
