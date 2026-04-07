package usecase

import (
	"log/slog"
	"tg-an/config"
	"tg-an/pkg/blog"
)

type Usecase struct {
	db       RepoPGInterface
	mdb      RepoMetricsPGInterface
	blclient *blog.Blog
	cfg      config.Analytics
	l        slog.Logger
}

func New(db RepoPGInterface, mdb RepoMetricsPGInterface, blclient *blog.Blog, cfg config.Analytics, l slog.Logger) *Usecase {
	return &Usecase{
		db:       db,
		mdb:      mdb,
		blclient: blclient,
		cfg:      cfg,
		l:        l,
	}
}
