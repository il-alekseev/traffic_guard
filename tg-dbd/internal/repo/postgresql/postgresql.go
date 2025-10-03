package postresql

import (
	"tg-dbd/pkg/logger"
	"tg-dbd/pkg/pgorm/pgorm"
)

type RepoPG struct {
	db     pgorm.Interface
	logger logger.Interface
}

func New(db pgorm.Interface, l logger.Interface) *RepoPG {
	return &RepoPG{
		db:     db,
		logger: l,
	}
}
