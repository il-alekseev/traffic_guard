package usecase

import "tg-dbd/pkg/logger"

type RepoPGInterface interface {
}

type Usecase struct {
	//db RepoPGInterface
	l logger.Interface
}

func New( /*db RepoPGInterface, */ l logger.Interface) *Usecase {
	return &Usecase{
		//db: db,
		l: l,
	}
}
