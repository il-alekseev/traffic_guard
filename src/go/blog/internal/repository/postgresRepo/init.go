package postgresRepo

import (
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB(db *gorm.DB) *PostgresDB {
	return &PostgresDB{
		db: db,
	}
}
