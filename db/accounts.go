package db

import (
	"context"
	"database/sql"
	"fmt"
)

type AccountsRepository struct {
	dbConn *sql.DB
}

type Account struct {
	ID       string
	UserName string
	Password string
}

func NewAccountsRepository(dbConn *sql.DB) AccountsRepository {
	return AccountsRepository{dbConn}
}

func (ac AccountsRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	row := ac.dbConn.QueryRowContext(
		ctx,
		"SELECT count(username) FROM accounts WHERE username=$1",
		username,
	)

	var count uint8

	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	if count >= 1 {
		return true, nil
	}

	return false, nil
}

func (ac AccountsRepository) Get(ctx context.Context, username string) (*Account, error) {
	row := ac.dbConn.QueryRowContext(
		ctx,
		"SELECT id, username, password FROM accounts WHERE username=$1",
		username,
	)

	var account Account

	err := row.Scan(&account.ID, &account.UserName, &account.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}

func (ac AccountsRepository) Save(ctx context.Context, username string, password string) (string, error) {
	stmt, err := ac.dbConn.PrepareContext(ctx, "INSERT INTO accounts (username, password) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var accountID string
	err = stmt.QueryRow(username, password).Scan(&accountID)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(accountID), nil
}
