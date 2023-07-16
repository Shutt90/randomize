package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type CockroachClient struct {
	ctx context.Context
	db  *sql.DB
}

type storedPassword struct {
	WebsiteName string    `json:"websiteName"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Created     time.Time `json:"created"`
}

func NewCockroachClient(ctx context.Context, db *sql.DB) *CockroachClient {
	return &CockroachClient{ctx: ctx, db: db}
}

func (cc *CockroachClient) Store(sp storedPassword) error {
	cc.db.ExecContext(cc.ctx, "CREATE TABLE IF NOT EXISTS password (websiteName varchar(255), username varchar(255), password varchar(255))")

	query := fmt.Sprintf("INSERT INTO password (websiteName, username, password) VALUES (%v, %v, %v);", &sp.WebsiteName, &sp.Username, &sp.Password)

	_, err := cc.db.ExecContext(cc.ctx, query)
	if err != nil {
		log.Fatal("failed to execute: ", err)
	}

	return nil
}

func (cc *CockroachClient) GetPassword(websiteName string) (string, error) {
	var password string
	err := cc.db.QueryRowContext(cc.ctx, "SELECT password FROM password WHERE websiteName = ?;", websiteName).Scan(password)
	if err != nil {
		return "", err
	}

	return password, nil
}

func (cc *CockroachClient) GetAllPasswords() ([]storedPassword, error) {
	rows, err := cc.db.QueryContext(cc.ctx, "SELECT websiteName, username, password FROM password;")
	if err != nil {
		if err == sql.ErrNoRows {
			return []storedPassword{}, fmt.Errorf("no passwords found")
		}
		return []storedPassword{}, fmt.Errorf("unknown error")
	}

	var passwords []storedPassword

	for rows.Next() {
		var password storedPassword
		err = rows.Scan(
			&password.WebsiteName,
			&password.Username,
			&password.Password,
		)
		if err != nil {
			fmt.Println(err)
			break
		}
		passwords = append(passwords, password)

	}

	return passwords, nil
}
