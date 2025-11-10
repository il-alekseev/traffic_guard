package pgorm

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Database - интерфейс для работы с базой данных
type Interface interface {
	// GetDB возвращает экземпляр *gorm.DB для работы с базой данных
	GetDB() *gorm.DB

	// Close закрывает соединение с базой данных
	Close() error

	// Migrate выполняет миграции для указанных моделей
	Migrate(models ...interface{}) error

	// WithTx выполняет переданную функцию в транзакции
	WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error

	// HealthCheck проверяет соединение с базой данных
	HealthCheck() error
}

// Проверка, что Postgres реализует интерфейс Database
var _ Interface = (*Postgres)(nil)

// Postgres - структура для работы с PostgreSQL через GORM
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	autoMigrate  bool
	models       []interface{}

	DB *gorm.DB
}

// New - создает новое подключение к PostgreSQL с использованием GORM
func New(dsn string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
		autoMigrate:  false,
		models:       make([]interface{}, 0),
	}

	// Применяем пользовательские опции
	for _, opt := range opts {
		opt(pg)
	}

	var err error
	for pg.connAttempts > 0 {
		pg.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	// Настраиваем пул соединений
	sqlDB, err := pg.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("postgres - failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(pg.maxPoolSize)
	sqlDB.SetMaxIdleConns(pg.maxPoolSize / 2)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Выполняем автоматические миграции, если включено
	if pg.autoMigrate && len(pg.models) > 0 {
		if err := pg.DB.AutoMigrate(pg.models...); err != nil {
			return nil, fmt.Errorf("postgres - failed to auto migrate: %w", err)
		}
	}

	return pg, nil
}

// Close - закрывает соединение с базой данных
func (p *Postgres) Close() error {
	if p.DB != nil {
		sqlDB, err := p.DB.DB()
		if err != nil {
			return fmt.Errorf("postgres - failed to get sql.DB for closing: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}

// Migrate - выполняет миграции для указанных моделей
func (p *Postgres) Migrate(models ...interface{}) error {
	fmt.Println("Running AutoMigrate for models:", models) // Логируем, какие модели передаются
	if err := p.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("postgres - failed to migrate: %w", err)
	}
	fmt.Println("Migration completed without errors") // Убедимся, что дошли до конца
	return nil
}

// WithTx выполняет функцию в транзакции
func (p *Postgres) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return p.DB.WithContext(ctx).Transaction(fn)
}

// GetDB возвращает экземпляр *gorm.DB для работы с базой данных
func (p *Postgres) GetDB() *gorm.DB {
	return p.DB
}

// HealthCheck проверяет соединение с базой данных
func (p *Postgres) HealthCheck() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	return sqlDB.Ping()
}
