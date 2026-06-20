package pgorm

import "time"

// Option - тип для функций опций конфигурации
type Option func(*Postgres)

// MaxPoolSize - устанавливает максимальный размер пула соединений
func MaxPoolSize(num int) Option {
	return func(p *Postgres) {
		if num != 0 {
			p.maxPoolSize = num
		}
	}
}

// ConnAttempts - устанавливает количество попыток подключения
func ConnAttempts(num int) Option {
	return func(p *Postgres) {
		if num != 0 {
			p.connAttempts = num
		}
	}
}

// ConnTimeout - устанавливает таймаут между попытками подключения
func ConnTimeout(timeout time.Duration) Option {
	return func(p *Postgres) {
		if timeout != 0 {
			p.connTimeout = timeout
		}
	}
}

// AutoMigrate - включает автоматическое создание/обновление таблиц
func AutoMigrate(enabled bool) Option {
	return func(c *Postgres) {
		c.autoMigrate = enabled
	}
}

// Models - добавляет модели для автоматического создания таблиц
func Models(models ...interface{}) Option {
	return func(c *Postgres) {
		c.models = append(c.models, models...)
	}
}
