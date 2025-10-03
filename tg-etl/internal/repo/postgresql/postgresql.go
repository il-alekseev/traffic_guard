package postgresql

import (
	"cmd/etl/pkg/logger"
	"cmd/etl/pkg/pgorm"
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
