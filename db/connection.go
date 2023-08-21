package db

import (
	"database/sql"

	"github.com/Maycon-Santos/go-snake-backend/process"
	_ "github.com/lib/pq"
)

func NewConnection(env *process.Env) (*sql.DB, error) {
	db, err := sql.Open(env.Database.Driver, env.Database.ConnURI)
	if err != nil {
		return nil, err
	}

	return db, nil
}
