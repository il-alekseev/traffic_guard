package usecase

import (
	"context"
	"fmt"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
)

func (uc *QueryUseCase) GetDeviceByID(ctx context.Context, id int) (*models.Device, error) {
	cacheKey := fmt.Sprintf("device:id:%d", id)
	// Пытаемся получить из кеша
	if cached, exists := uc.c.Get(cacheKey); exists {
		//uc.l.DebugContext(ctx, "cache hit for device", wsl.Int("id", id))
		return cached.(*models.Device), nil
	}
	//uc.l.DebugContext(ctx, "cache miss for device", slog.Int("id", id))

	// Если нет в кеше, ищем в БД
	device, err := uc.etlDB.GetDeviceByID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get device by id from db", wsl.Err((err)))
	}
	// Если устройство есть в базе, добавляем в кеш
	if device != nil {
		//uc.l.DebugContext(ctx, "success got device by id", slog.Int("id", id))
		uc.c.Set(cacheKey, device)
	}
	return device, nil
}

func (uc *QueryUseCase) CreateDevice(ctx context.Context, device models.Device) error {
	err := uc.etlDB.CreateDevice(ctx, device)
	if err != nil {
		uc.l.ErrorContext(ctx, "create device", wsl.Err((err)))
		return err
	}
	// Инвалидируем возможные кеши
	uc.c.Delete(fmt.Sprintf("device:id:%d", device.ID))
	//uc.l.DebugContext(ctx, "success create device")
	return nil
}

func (uc *QueryUseCase) GetDevices(ctx context.Context) ([]models.Device, error) {
	cacheKey := "devices:all"
	if cached, exists := uc.c.Get(cacheKey); exists {
		//uc.l.DebugContext(ctx, "cache hit for all devices")
		return cached.([]models.Device), nil
	}
	//uc.l.DebugContext(ctx, "cache miss for all devices")
	devices, err := uc.etlDB.GetDevices(ctx)
	if err != nil {
		uc.l.ErrorContext(ctx, "get devices", wsl.Err((err)))
		return nil, err
	}
	uc.c.Set(cacheKey, devices)
	//uc.l.DebugContext(ctx, "success got devices")
	return devices, nil
}
