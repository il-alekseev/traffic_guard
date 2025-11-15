package usecase

import (
	"log/slog"
	"tg-an/pkg/blog"
)

type Usecase struct {
	db       RepoPGInterface
	mdb      RepoMetricsPGInterface
	blclient *blog.Blog
	l        slog.Logger
}

func New(db RepoPGInterface, mdb RepoMetricsPGInterface, blclient *blog.Blog, l slog.Logger) *Usecase {
	return &Usecase{
		db:       db,
		mdb:      mdb,
		blclient: blclient,
		l:        l,
	}
}
