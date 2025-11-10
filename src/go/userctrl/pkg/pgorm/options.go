package pgorm

import "time"

// Option - тип для функций опций конфигурации
type Option func(*Postgres)

// MaxPoolSize - устанавливает максимальный размер пула соединений
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// ConnAttempts - устанавливает количество попыток подключения
func ConnAttempts(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

// ConnTimeout - устанавливает таймаут между попытками подключения
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
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
