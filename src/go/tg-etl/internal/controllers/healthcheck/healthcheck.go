package healthcheck

import (
	"cmd/etl/pkg/pgorm"
	"cmd/etl/pkg/slogger/wsl"
	"context"
	"fmt"
	"log/slog"
	"time"
)

type HealthCheckService struct {
	db       pgorm.Interface
	l        slog.Logger
	interval time.Duration
}

func NewHealthCheckService(db pgorm.Interface, l slog.Logger, interval time.Duration) *HealthCheckService {
	return &HealthCheckService{
		db:       db,
		l:        l,
		interval: interval,
	}
}

// HealthStatus представляет статус здоровья системы
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  bool      `json:"database"`
	Cache     bool      `json:"cache"`
	Errors    []string  `json:"errors,omitempty"`
}

// Start запускает периодическую проверку здоровья
func (s *HealthCheckService) Start(ctx context.Context) {
	s.l.InfoContext(ctx, "starting health check service", slog.Duration("interval", s.interval))

	// Первая проверка сразу при старте
	if err := s.checkHealth(ctx); err != nil {
		s.l.ErrorContext(ctx, "initial health check failed", wsl.Err(err))
	} else {
		s.l.InfoContext(ctx, "initial health check passed")
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.l.InfoContext(ctx, "stopping health check service")
			return
		case <-ticker.C:
			if err := s.checkHealth(ctx); err != nil {
				s.l.ErrorContext(ctx, "health check failed", wsl.Err(err))
			} else {
				s.l.DebugContext(ctx, "health check passed")
			}
		}
	}
}

// checkHealth выполняет комплексную проверку здоровья системы
func (s *HealthCheckService) checkHealth(ctx context.Context) error {
	healthStatus := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Database:  true,
		Cache:     true, // Для in-memory кеша всегда true
	}

	// Проверка базы данных
	if err := s.checkDatabase(ctx); err != nil {
		healthStatus.Database = false
		healthStatus.Status = "degraded"
		healthStatus.Errors = append(healthStatus.Errors, fmt.Sprintf("database: %v", err))
	}

	// Проверка основных таблиц
	if healthStatus.Database {
		if err := s.checkTables(ctx); err != nil {
			healthStatus.Status = "degraded"
			healthStatus.Errors = append(healthStatus.Errors, fmt.Sprintf("tables: %v", err))
		}
	}

	// Проверка производительности
	if healthStatus.Database {
		if err := s.checkPerformance(ctx); err != nil {
			healthStatus.Status = "degraded"
			healthStatus.Errors = append(healthStatus.Errors, fmt.Sprintf("performance: %v", err))
		}
	}

	// Логируем результат проверки
	s.logHealthStatus(ctx, healthStatus)

	if healthStatus.Status != "healthy" {
		return fmt.Errorf("health check failed: %v", healthStatus.Errors)
	}

	return nil
}

// checkDatabase проверяет подключение к базе данных
func (s *HealthCheckService) checkDatabase(ctx context.Context) error {
	// Получаем native sql.DB
	sqlDB, err := s.db.GetDB().DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// Проверяем подключение
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Проверяем статистику подключений
	stats := sqlDB.Stats()
	if stats.OpenConnections > stats.MaxOpenConnections*80/100 {
		s.l.WarnContext(ctx, "database connections approaching limit",
			slog.Int("open_connections", stats.OpenConnections),
			slog.Int("max_connections", stats.MaxOpenConnections))
	}

	return nil
}

// checkTables проверяет наличие и доступность основных таблиц
func (s *HealthCheckService) checkTables(ctx context.Context) error {
	tables := []string{"ids_logs", "sessions", "sources", "domains", "devices"}

	for _, table := range tables {
		if err := s.checkTableExists(ctx, table); err != nil {
			return fmt.Errorf("table %s check failed: %w", table, err)
		}
	}

	return nil
}

// checkTableExists проверяет существование таблицы
func (s *HealthCheckService) checkTableExists(ctx context.Context, tableName string) error {
	var exists bool
	err := s.db.GetDB().WithContext(ctx).
		Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = ?)", tableName).
		Scan(&exists).Error

	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	// Проверяем что можем выполнить простой запрос
	var count int64
	err = s.db.GetDB().WithContext(ctx).Table(tableName).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to count rows in %s: %w", tableName, err)
	}

	s.l.DebugContext(ctx, "table check passed",
		slog.String("table", tableName),
		slog.Int64("row_count", count))

	return nil
}

// checkPerformance проверяет производительность базы данных
func (s *HealthCheckService) checkPerformance(ctx context.Context) error {
	// Проверяем время выполнения простого запроса
	start := time.Now()

	var result int
	err := s.db.GetDB().WithContext(ctx).Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		return fmt.Errorf("performance check query failed: %w", err)
	}

	duration := time.Since(start)
	if duration > 100*time.Millisecond {
		s.l.WarnContext(ctx, "database response time is high",
			slog.Duration("response_time", duration))
	}

	// Проверяем размеры таблиц
	if err := s.checkTableSizes(ctx); err != nil {
		s.l.WarnContext(ctx, "failed to check table sizes", wsl.Err(err))
	}

	return nil
}

// checkTableSizes проверяет размеры таблиц и предупреждает о больших таблицах
func (s *HealthCheckService) checkTableSizes(ctx context.Context) error {
	type TableSize struct {
		TableName string `gorm:"column:table_name"`
		Size      string `gorm:"column:table_size"`
	}

	var tableSizes []TableSize
	err := s.db.GetDB().WithContext(ctx).Raw(`
		SELECT 
			table_name,
			pg_size_pretty(pg_total_relation_size(quote_ident(table_name))) as table_size
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY pg_total_relation_size(quote_ident(table_name)) DESC
	`).Scan(&tableSizes).Error

	if err != nil {
		return err
	}

	for _, ts := range tableSizes {
		s.l.DebugContext(ctx, "table size",
			slog.String("table", ts.TableName),
			slog.String("size", ts.Size))
	}

	return nil
}

// logHealthStatus логирует результат проверки здоровья
func (s *HealthCheckService) logHealthStatus(ctx context.Context, status *HealthStatus) {
	logger := s.l.With(
		slog.String("status", status.Status),
		slog.Bool("database", status.Database),
		slog.Bool("cache", status.Cache),
		slog.Time("timestamp", status.Timestamp),
	)

	if status.Status == "healthy" {
		logger.InfoContext(ctx, "health check passed")
	} else {
		logger.ErrorContext(ctx, "health check failed",
			slog.Any("errors", status.Errors))
	}
}

// GetDetailedStatus возвращает детальный статус для API
func (s *HealthCheckService) GetDetailedStatus(ctx context.Context) (*HealthStatus, error) {
	healthStatus := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Database:  true,
		Cache:     true,
	}

	// Проверка базы данных
	if err := s.checkDatabase(ctx); err != nil {
		healthStatus.Database = false
		healthStatus.Status = "degraded"
		healthStatus.Errors = append(healthStatus.Errors, fmt.Sprintf("database: %v", err))
	}

	// Дополнительные проверки для детального статуса
	if healthStatus.Database {
		// Проверяем количество записей в основных таблицах
		if counts, err := s.getTableCounts(ctx); err != nil {
			healthStatus.Errors = append(healthStatus.Errors, fmt.Sprintf("counts: %v", err))
		} else {
			// Можно добавить counts в healthStatus если нужно
			s.l.DebugContext(ctx, "table counts", slog.Any("counts", counts))
		}

		// Проверяем последнюю обработанную запись
		if lastLog, err := s.getLastProcessedLog(ctx); err != nil {
			healthStatus.Errors = append(healthStatus.Errors, fmt.Sprintf("last_log: %v", err))
		} else if lastLog != nil {
			s.l.DebugContext(ctx, "last processed log",
				slog.Uint64("log_id", uint64(lastLog.ID)),
				slog.Time("timestamp", lastLog.Timestamp))
		}
	}

	return healthStatus, nil
}

// getTableCounts возвращает количество записей в основных таблицах
func (s *HealthCheckService) getTableCounts(ctx context.Context) (map[string]int64, error) {
	counts := make(map[string]int64)
	tables := []string{"ids_logs", "sessions", "sources", "domains", "devices"}

	for _, table := range tables {
		var count int64
		err := s.db.GetDB().WithContext(ctx).Table(table).Count(&count).Error
		if err != nil {
			return nil, fmt.Errorf("failed to count %s: %w", table, err)
		}
		counts[table] = count
	}

	return counts, nil
}

// getLastProcessedLog возвращает последнюю запись из IdsLogs
func (s *HealthCheckService) getLastProcessedLog(ctx context.Context) (*struct {
	ID        uint      `gorm:"column:id"`
	Timestamp time.Time `gorm:"column:timestamp"`
}, error) {
	var result struct {
		ID        uint      `gorm:"column:id"`
		Timestamp time.Time `gorm:"column:timestamp"`
	}

	err := s.db.GetDB().WithContext(ctx).
		Table("ids_logs").
		Select("id, timestamp").
		Order("id DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CheckReady проверяет готовность сервиса (для Kubernetes readiness probe)
func (s *HealthCheckService) CheckReady(ctx context.Context) bool {
	status, err := s.GetDetailedStatus(ctx)
	if err != nil {
		s.l.ErrorContext(ctx, "readiness check failed", wsl.Err(err))
		return false
	}
	return status.Status == "healthy"
}

// CheckLive проверяет живучесть сервиса (для Kubernetes liveness probe)
func (s *HealthCheckService) CheckLive(ctx context.Context) bool {
	// Более легкая проверка чем CheckReady
	if err := s.checkDatabase(ctx); err != nil {
		s.l.ErrorContext(ctx, "liveness check failed", wsl.Err(err))
		return false
	}
	return true
}
