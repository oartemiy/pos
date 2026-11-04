package jsonfile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// json storage for T object
type store[T any] struct {
	mutex sync.Mutex
	path  string
}

func newStore[T any](path string) *store[T] {
	return &store[T]{path: path}
}

func (s *store[T]) load() ([]T, error) {
	file, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("jsonfile: read %s: %w", s.path, err)
	}
	var items []T
	if err := json.Unmarshal(file, &items); err != nil {
		return nil, fmt.Errorf("jsonfile: decode %s: %w", s.path, err)
	}
	return items, nil
}

func (s *store[T]) read(ctx context.Context) ([]T, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.load()
}

func (s *store[T]) update(ctx context.Context, updateFn func([]T) ([]T, error)) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	items, err := s.load()
	if err != nil {
		return err
	}
	items, err = updateFn(items)
	if err != nil {
		return err
	}
	return s.save(items)
}

func (s *store[T]) save(items []T) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("jsonfile: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return fmt.Errorf("jsonfile: encode: %w", err)
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("jsonfile: write: %w", err)
	}
	return os.Rename(tmp, s.path)
}
