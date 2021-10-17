package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const DriverPostgres = "postgres"

type connectionProvider interface {
	GetDriverName() string
	GetUsername() string
	GetPassword() string
	GetHost() string
	GetPort() string
	GetDatabase() string
	GetSchema() string
}

func NewPostgresConnection(p connectionProvider) (*sqlx.DB, error) {
	// Disable ssl (for now)
	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s search_path=%s sslmode=disable",
		p.GetUsername(),
		p.GetPassword(),
		p.GetDatabase(),
		p.GetHost(),
		p.GetPort(),
		p.GetSchema(),
	)

	return sqlx.Connect("postgres", connString)
}
