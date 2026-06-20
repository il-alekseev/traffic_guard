package postgresql

import (
	"log/slog"
	"tg-etl/pkg/pgorm"
)

type ELTRepoPG struct {
	db pgorm.Interface
	l  slog.Logger
}

func NewELTRepoPG(db pgorm.Interface, l slog.Logger) *ELTRepoPG {
	return &ELTRepoPG{
		db: db,
		l:  l,
	}
}
