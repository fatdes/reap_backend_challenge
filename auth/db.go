package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

type DBPGX struct {
	URL string
}

func (db *DBPGX) Exists(username string) (bool, error) {
	conn, err := pgx.Connect(context.Background(), db.URL)
	if err != nil {
		return false, err
	}
	defer conn.Close(context.Background())

	var exists bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM auth WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *DBPGX) Create(username string, password string) error {
	conn, err := pgx.Connect(context.Background(), db.URL)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "INSERT INTO auth (username, password, created_at) VALUES ($1, $2, $3)", username, password, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (db *DBPGX) Login(username string, password string) (bool, error) {
	conn, err := pgx.Connect(context.Background(), db.URL)
	if err != nil {
		return false, err
	}
	defer conn.Close(context.Background())

	var succeed bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM auth WHERE username = $1 AND password = $2)", username, password).Scan(&succeed)
	if err != nil {
		return false, err
	}
	return succeed, nil
}
