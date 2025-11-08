package geo

import (
	"context"
	"errors"
	"net"

	"github.com/oschwald/maxminddb-golang"

	"scrapper/internal/models"
)

type MaxMindProvider struct {
	db *maxminddb.Reader
}

func NewMaxMindProvider(dbPath string) (*MaxMindProvider, error) {
	reader, err := maxminddb.Open(dbPath)
	if err != nil {
		return nil, err
	}

	return &MaxMindProvider{db: reader}, nil
}

func (p *MaxMindProvider) Lookup(ctx context.Context, ip net.IP) (models.GeoInfo, error) {
	_ = ctx

	if p == nil || p.db == nil {
		return nil, errors.New("maxmind provider is not initialized")
	}

	if ip == nil {
		return nil, errors.New("ip is nil")
	}

	var record map[string]any
	if err := p.db.Lookup(ip, &record); err != nil {
		return nil, err
	}

	if len(record) == 0 {
		return nil, nil
	}

	return models.GeoInfo(record), nil
}

func (p *MaxMindProvider) Close() error {
	if p == nil || p.db == nil {
		return nil
	}

	return p.db.Close()
}
