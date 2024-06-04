package main

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}
}
