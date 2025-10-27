package usecase

type UseCase struct {
}

// New — конструктор UseCase, принимающий конфигурацию и зависимости, необходимые для работы бизнес-логики
func New() *UseCase {
	uc := UseCase{}
	return &uc
}
