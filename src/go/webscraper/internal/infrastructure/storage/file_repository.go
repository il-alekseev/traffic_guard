package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"scrapper/internal/models"
)

type FileRepository struct {
	path string
	mu   sync.Mutex
}

func NewFileRepository(path string) (*FileRepository, error) {
	if path == "" {
		return nil, errors.New("file repository path is empty")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	return &FileRepository{path: path}, nil
}

func (r *FileRepository) SaveAttempt(ctx context.Context, attempt models.ContentAttempt) error {
	_ = ctx

	if attempt.URL == "" {
		return errors.New("attempt url is empty")
	}

	payload, err := json.Marshal(attempt)
	if err != nil {
		return err
	}

	payload = append(payload, '\n')

	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile(r.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(payload)
	return err
}
