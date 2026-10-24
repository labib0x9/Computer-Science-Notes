package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // without import, gives error
)

// CREATE DATABASE testDb;
// SHOW DATABASES;
// USE testDb;

/**
CREATE TABLE items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL
);
**/

var (
	user     = "root"
	password = ""
	ip       = "127.0.0.1"
	port     = "3306"
)

func main() {

	dbName := "testDb"
	dsn := user + ":" + password + "@tcp(" + ip + ":" + port + ")/" + dbName
	fmt.Println(dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Open error : %w", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("Ping error")
		return
	}

	// Connected to mysql
	fmt.Println("Connected")

	insertStmt := `INSERT INTO items (content) VALUES (?)` // ? is replaceable with variable
	queryStmt := `SELECT content FROM items WHERE id = ?`

	// Insert
	content := "Word"
	res, err := db.Exec(insertStmt, content)
	if err != nil {
		return
	}
	res.LastInsertId()
	res.RowsAffected()

	// Query
	var storedContent string
	index := 1
	row := db.QueryRow(queryStmt, index)
	if err := row.Scan(&storedContent); err != nil {
		return
	}

	// Get last insert ID (simulate auto increment)
	var lastID int
	row = db.QueryRow(`SELECT MAX(id) FROM items`)
	err = row.Scan(&lastID)
}

func createDb(dbName string) error {
	dsn := user + ":" + password + "@tcp(" + ip + ":" + port + ")/"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return err
	}

	// database created
	return nil
}
