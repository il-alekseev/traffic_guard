package usecase

import (
	"log/slog"
)

type Usecase struct {
	db  RepoPGInterface
	mdb RepoMetricsPGInterface
	l   slog.Logger
}

func New(db RepoPGInterface, mdb RepoMetricsPGInterface, l slog.Logger) *Usecase {
	return &Usecase{
		db:  db,
		mdb: mdb,
		l:   l,
	}
}
