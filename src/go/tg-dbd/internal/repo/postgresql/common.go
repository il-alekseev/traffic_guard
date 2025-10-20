package postresql

import (
	"context"
	"fmt"
)

func (r *RepoPG) GetDevices(ctx context.Context) ([]string, error) {
	var devices []string
	query := r.db.GetDB().WithContext(ctx).Table("devices").
		Select(`host_name`)

	if err := query.Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	return devices, nil
}
