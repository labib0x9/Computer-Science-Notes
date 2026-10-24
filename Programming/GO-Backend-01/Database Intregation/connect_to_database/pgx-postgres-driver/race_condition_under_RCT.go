package main

/*
	pgx  -  POSTGRES DRIVER

	Race condition under Read committed isolation level...
	// default isolation level is read committed.
*/

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

//// ---- Change isolation level ---- ////
/*
tx, err := pool.BeginTx(context.Background(), pgx.TxOptions{
	IsoLevel: pgx.RepeatableRead,
})
*/

// Read Committed isolation level transaction...
func txRC(pool *pgxpool.Pool, id int, serial int) (ans any, err error) {

	tx, err := pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background()) // safe - after commit rollback does nothing

	query := "select * from urls where id=$1" // row - lock here to prevent...
	row := tx.QueryRow(context.Background(), query, id)
	var found Url
	if err := row.Scan(&found.Id, &found.Url); err != nil {
		return nil, err
	}

	if found.Url != "https://labib0x0hunter.netlify.app/#projects" {
		return nil, errors.New("not the url i want...")
	}

	longUrl := "https://labib0x0hunter.netlify.app/no_need_to_know_whoami_" + strconv.Itoa(serial)
	query = "update urls set url=$2 where id=$1 returning id, url"
	updatedRow := tx.QueryRow(context.Background(), query, id, longUrl)
	var updatedUrl Url
	if err := updatedRow.Scan(&updatedUrl.Id, &updatedUrl.Url); err != nil {
		return nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return updatedUrl, nil
}

func main() {
	//// ---- Connection Pool ---- ////
	// pgURL := "postgres://username:password@localhost:5432/database_name"
	pgURL := "postgres://short:secret@localhost:5432/short"
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

	/// ---- ///
	var wg sync.WaitGroup
	for i := range 4 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			result, err := txRC(pool, 6, id)
			fmt.Println("GOROUTINE=", id, "REsult=", result, "Error=", err)
		}(i)
	}

	wg.Wait()
}
