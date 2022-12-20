package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type storedPassword struct {
	WebsiteName string    `json:"websiteName"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Created     time.Time `json:"created"`
}

func main() {
	fmt.Println("Enter website name: ")
	var entry storedPassword
	_, err := fmt.Scanln(&entry.WebsiteName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Enter username: ")
	_, err = fmt.Scanln(&entry.Username)
	if err != nil {
		fmt.Println(err)
		return
	}

	// scanner.Scan(entry.WebsiteName)
	fmt.Println("Number of characters")
	var chars string
	n, err := fmt.Scanln(&chars)
	if err != nil || n > 3 {
		fmt.Println("Error: can't be higher than 255 and ", err)
		return
	}
	convertedNum, err := strconv.Atoi(chars)
	if err != nil {
		fmt.Println(err)
		return
	}
	if convertedNum > 255 || convertedNum < 8 {
		fmt.Println("Error: can't be higher than 255 or lower than 8")
		return
	}

	entry.Password = randomize(uint8(convertedNum))

	entry.store()
}

func randomize(numLetters uint8) string {
	rand.Seed(time.Now().UnixNano())
	var password []string
	var i uint8
	for i = 0; i < numLetters; i++ {
		rand := rand.Intn(126-33) + 33
		letter := fmt.Sprintf("%c", rand)
		password = append(password, string(letter))
	}

	var sep string
	joinedPassword := strings.Join(password, sep)

	return joinedPassword
}

func (sp *storedPassword) store() error {
	godotenv.Load()
	dsn := os.Getenv("DB_DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		fmt.Println("could not connect to database")
		return nil
	}

	conn.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS password (websiteName varchar(255), username varchar(255), password varchar(255))")

	query := fmt.Sprintf("INSERT INTO password (websiteName, username, password) VALUES (%v, %v, %v);", &sp.WebsiteName, &sp.Username, &sp.Password)

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		log.Fatal("failed to execute query", err)
	}

	fmt.Println("successfully stored password for ", sp.WebsiteName)

	return nil
}
