package usecase

import (
	"log/slog"
)

type Usecase struct {
	db RepoPGInterface
	l  slog.Logger
}

func New(db RepoPGInterface, l slog.Logger) *Usecase {
	return &Usecase{
		db: db,
		l:  l,
	}
}
