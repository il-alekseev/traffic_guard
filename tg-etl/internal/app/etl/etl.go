package etl

import (
	"cmd/etl/config"
	"cmd/etl/internal/models"
	"cmd/etl/internal/repo/postgresql"
	"cmd/etl/pkg/logger"
	"cmd/etl/pkg/pgorm"
	"fmt"
)

func getDSN(host, user, pass, dbname, port, sslmode string) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, pass, dbname, port, sslmode,
	)
}

func createConnection(dsn string, maxPoolSize int, l *logger.Logger, models ...interface{}) (*postgresql.RepoPG, error) {
	db, err := pgorm.New(
		dsn,
		pgorm.MaxPoolSize(maxPoolSize),
		pgorm.AutoMigrate(true),
		pgorm.Models(models...),
	)
	if err != nil {
		l.Error("Error with creation connetcion %v", err)
		return nil, err
	}
	l.Debug("connection created successfully")
	if err := db.HealthCheck(); err != nil {
		l.Error("database health check failed: %v", err)
		return nil, err
	} else {
		l.Info("database connection is healthy")
	}
	pg := postgresql.New(db, l)
	return pg, nil
}

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	dsn := getDSN(cfg.PG.Host, cfg.PG.User, cfg.PG.Pass,
		cfg.PG.DBName, cfg.PG.Port, cfg.PG.SSLMode)
	db, err := createConnection(dsn, cfg.PG.PoolMax, l,
		// Автомиграция только базовых таблиц
		models.Category{},
		models.CategoryDomain{},
		models.Decision{},
		models.Detection{},
		models.Device{},
		models.Domain{},
		models.Session{},
		models.Source{},
		models.Status{},
		// В release убрать!
		models.IdsLog{},
	)

	if err != nil {
		l.Fatal("Create db connection error: %v", err)
		return
	}
	l.Debug("db info %v", db)

	// Заполняем структуру тестовыми данными
	//logs, err := parser.ParseCSVFileToIdsLog("C:/Users/Asus/Projects/VS Code/Continent/traffic_guard/src/go/tg-etl/test_data/security_log-2025-9-20_15-58.csv")
	//if err != nil {
	//	l.Fatal("could not parse csv file: %v", err)
	//	return
	//}
	//l.Debug("got %d logs", len(logs))
	//ctx := context.Background()
	//for _, log := range logs {
	//	if err := db.CreateLog(&ctx, &log); err != nil {
	//		fmt.Printf("failed to create log: %s\n", err)
	//	}
	//}
}
