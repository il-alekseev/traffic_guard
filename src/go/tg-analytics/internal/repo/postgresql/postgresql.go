package postresql

import (
	"log/slog"
	"tg-an/pkg/pgorm/pgorm"
)

type RepoPG struct {
	db pgorm.Interface
	l  slog.Logger
}

func New(db pgorm.Interface, l slog.Logger) *RepoPG {
	return &RepoPG{
		db: db,
		l:  l,
	}
}
