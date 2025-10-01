package postgresGorm

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log/slog"
	"time"

	"fiermon-blog/config"
	"fiermon-blog/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type PostgresDB struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	db           *gorm.DB
}

// NewPostgres - подключение к БД Postgres с реализацией интерфейса Read
func NewPostgres(cfg config.PG, log *slog.Logger) (*PostgresDB, error) {
	pl := "Repo"
	log = log.With(
		slog.String("layer", pl),
	)

	log.Info("Connecting to postgres database")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
	)

	pg := &PostgresDB{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	var err error
	for pg.connAttempts > 0 {
		pg.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			log.Info(fmt.Sprintf("Connected to postgres database, on attempt: %d", _defaultConnAttempts-pg.connAttempts))
			break
		}

		log.Info(fmt.Sprintf("Postgres is trying to connect, attempts left: %d", pg.connAttempts))
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	// Настраиваем пул соединений
	sqlDB, err := pg.db.DB()
	if err != nil {
		return nil, fmt.Errorf("postgres - failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(pg.maxPoolSize)
	sqlDB.SetMaxIdleConns(pg.maxPoolSize / 2)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Настройка автомиграции
	err = pg.db.AutoMigrate(&models.BusinessLog{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate models: %w", err)
	}

	return pg, nil
}
