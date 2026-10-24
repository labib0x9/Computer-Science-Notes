package main

/*
	pgx  -  POSTGRES DRIVER
*/

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

/*
go get "github.com/jackc/pgx/v5"

// Table
CREATE TABLE urls (
     id  SERIAL PRIMARY KEY,
     url TEXT         NOT NULL
);

*/

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

	//// ---- Configs ---- ////
	cfg := conn.Config()
	fmt.Println("DB=", cfg.Database)
	fmt.Println("Password=", cfg.Password)

	/// --- QUERY --- ///

	//// ---- Insert / exec ---- ////
	insertQuery := "insert into urls(url) values($1)"
	cmds, err := conn.Exec(context.Background(), insertQuery, "https://labib0x0hunter.netlify.app/#projects")
	if err != nil {
		panic(err)
	}
	// cmds= INSERT 0 1
	fmt.Println("cmds=", cmds)

	//// ---- Multiple Rows ---- ////
	query := "select id, url from urls"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var url Url
	for rows.Next() {
		if err := rows.Scan(&url.Id, &url.Url); err != nil {
			panic(err)
		}
		fmt.Println("url=", url)
	}

	urls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Url])
	if err != nil {
		panic(err)
	}
	fmt.Println("urls=", urls)

	//// --- Update --- ////
	updateQuery := "update urls set url = $1 where id = $2"
	cmds, err = conn.Exec(context.Background(), updateQuery, "https://labib0x0hunter.netlify.app/whoami", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println("cmds=", cmds)
	fmt.Println("AffectedRows=", cmds.RowsAffected())

	updateQuery = "update urls set url = $1 where id = $2 returning *"
	row := conn.QueryRow(context.Background(), updateQuery, "https://labib0x0hunter.netlify.app/againUpdated", 1)
	var updatedUrl Url
	if err := row.Scan(&updatedUrl.Id, &updatedUrl.Url); err != nil {
		panic(err)
	}
	fmt.Println("updatedUrl=", updatedUrl)

	/// -- Single Row -- ////
	query = "select url from urls where id = $1"
	row = conn.QueryRow(context.Background(), query, 1)
	var checkUrl string
	if err = row.Scan(&checkUrl); err != nil {
		panic(err)
	}
	fmt.Println("(I should return first) checkUrl=", checkUrl)

	/// -- Count Rows -- ///
	countQuery := "select count(*) from urls"
	countRow := conn.QueryRow(context.Background(), countQuery)
	var count int
	if err = countRow.Scan(&count); err != nil {
		panic(err)
	}
	fmt.Println("count=", count)

	/// --- TRANSACTION --- ///
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted, // isolation level
	})
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec(context.Background(), "insert into urls(url) values($1)", "https://labib0x0hunter.netlify.app/#tx1")
	if err != nil {
		tx.Rollback(context.Background())
		panic(err)
	}

	_, err = tx.Exec(context.Background(), "insert into urls(url) values($1)", "https://labib0x0hunter.netlify.app/#tx2")
	if err != nil {
		tx.Rollback(context.Background())
		panic(err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		panic(err)
	}
}
