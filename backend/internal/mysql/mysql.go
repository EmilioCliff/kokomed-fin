package mysql

import "database/sql"

type Store struct {
	db *sql.DB
}

func NewStore() *Store {
	return &Store{}
}
