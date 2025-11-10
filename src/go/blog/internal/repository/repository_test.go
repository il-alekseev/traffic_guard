package repository

import (
	"context"
	"fiermon-blog/config"
	"fiermon-blog/internal/models"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"fiermon-blog/internal/repository/postgresGorm"
)

type Read interface {
	GetRecord(ctx context.Context, userMeta *models.UserMeta, page, limit int, role, contextID, search string) ([]models.BusinessLog, models.Meta, error)
}

func TestRepository(t *testing.T) {
	configPG := config.PG{
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "postgresPass",
		DBName:   "blog",
		SSLMode:  "disable",
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	read, err := postgresGorm.NewPostgres(configPG, log)
	if err != nil {
		t.Error("readWriter error:", err)
		return
	}

	type args struct {
		ctx       context.Context
		userMeta  *models.UserMeta
		page      int
		limit     int
		role      string
		contextID string
		search    string
	}

	tests := []struct {
		name    string
		fields  Read
		args    args
		want    []models.BusinessLog
		want1   models.Meta
		wantErr bool
	}{
		{
			name:   "Success",
			fields: read,
			args: args{
				ctx: context.Background(),
				userMeta: &models.UserMeta{
					ContextID:  "qwe",
					Username:   "admin",
					ClientRole: "CA",
				},
				page:      1,
				limit:     10,
				role:      "admin",
				contextID: "qwe",
				search:    "",
			},
			want:    []models.BusinessLog{},
			want1:   models.Meta{Limit: 10, Page: 1, Total: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1, err := tt.fields.GetRecord(tt.args.ctx, tt.args.userMeta, tt.args.page, tt.args.limit, tt.args.role, tt.args.contextID, tt.args.search)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRecord() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetRecord() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

}
