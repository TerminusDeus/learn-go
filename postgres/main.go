package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Book struct {
	isbn   string
	title  string
	author string
	price  float32
}

func main() {
	user := "postgres"
	pass := "postgres"
	host := "localhost"
	dbName := "mydb"

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, host, dbName))
	isErr(err)

	rows, err := db.Query("SELECT * FROM books")
	isErr(err)
	defer rows.Close()

	bks := make([]*Book, 0)
	for rows.Next() {
		bk := new(Book)
		err := rows.Scan(&bk.isbn, &bk.title, &bk.author, &bk.price)
		isErr(err)
		bks = append(bks, bk)
	}
	isErr(rows.Err())

	for _, bk := range bks {
		fmt.Printf("%s, %s, %s, £%.2f\n", bk.isbn, bk.title, bk.author, bk.price)
	}
}

func isErr(e error) {
	if e != nil {
		panic(e)
	}
}