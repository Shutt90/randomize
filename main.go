package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	cockroachDB "github.com/shutt90/password-generator/db"
	gui "github.com/shutt90/password-generator/gui"
)

func main() {
	err := godotenv.Load()
	ctx := context.Background()
	if err != nil {
		panic(err)
	}

	dsn := os.Getenv("DB_DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	defer db.Close()

	conn, err := db.Conn(ctx)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	cc := cockroachDB.NewCockroachClient(ctx, db)

	passwords, err := cc.GetAllPasswords()
	if err != nil {
		panic(err)
	}

	gui.MainWindow(cc, passwords)

}
