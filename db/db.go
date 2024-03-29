package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrCouldNotStore = errors.New("Unable to store password at this time")
	ErrBlankInput    = errors.New("One or more fields is empty")
)

type CockroachClient struct {
	ctx context.Context
	db  *sql.DB
}

type StoredPassword struct {
	WebsiteName string    `json:"websiteName"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Created     time.Time `json:"created"`
}

func NewCockroachClient(ctx context.Context, db *sql.DB) *CockroachClient {
	return &CockroachClient{ctx: ctx, db: db}
}

func (cc *CockroachClient) Store(sp StoredPassword) error {
	if sp.Password == "" || sp.Username == "" || sp.WebsiteName == "" {
		return ErrBlankInput
	}

	query := "INSERT INTO password (websiteName, username, password) VALUES ( $1, $2, $3 )"

	_, err := cc.db.ExecContext(cc.ctx, query, sp.WebsiteName, sp.Username, sp.Password)
	if err != nil {
		return ErrCouldNotStore
	}

	return nil
}

func (cc *CockroachClient) GetAllPasswords() ([]StoredPassword, error) {
	rows, err := cc.db.QueryContext(cc.ctx, "SELECT websiteName, username, password FROM password;")
	if err != nil {
		if err == sql.ErrNoRows {
			return []StoredPassword{}, fmt.Errorf("no passwords found")
		}
		return []StoredPassword{}, fmt.Errorf("unknown error")
	}

	var passwords []StoredPassword

	for rows.Next() {
		var password StoredPassword
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
