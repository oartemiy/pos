package jsonfile

import (
	"context"
	"pos/internal/domain/shared"
	"pos/internal/domain/student"
)

type StudentRepo struct {
	store *store[studentRecord]
}

// To check intercace implementation
var _ student.Repository = (*StudentRepo)(nil)

func NewStudentRepo(path string) *StudentRepo {
	return &StudentRepo{store: newStore[studentRecord](path)}
}

func (r *StudentRepo) Save(ctx context.Context, s *student.Student) error {
	record := toStudentRecord(s)
	return r.store.update(ctx, func(items []studentRecord) ([]studentRecord, error) {
		for i, item := range items {
			if item.ID == record.ID {
				items[i] = record
				return items, nil
			}
		}
		return append(items, record), nil
	})
}

func (r *StudentRepo) FindByID(ctx context.Context, id shared.ID) (*student.Student, error) {
	return r.findOnce(ctx, func(rec studentRecord) bool {
		return rec.ID == id.String()
	})
}

func (r *StudentRepo) FindByEmail(ctx context.Context, email shared.Email) (*student.Student, error) {
	return r.findOnce(ctx, func(rec studentRecord) bool {
		return rec.Email == email.String()
	})
}

func (r *StudentRepo) List(ctx context.Context) ([]*student.Student, error) {
	items, err := r.store.load()
	if err != nil {
		return nil, err
	}
	result := make([]*student.Student, 0, len(items))
	for _, rec := range items {
		s, err := toStudent(rec)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *StudentRepo) Delete(ctx context.Context, id shared.ID) error {
	return r.store.update(ctx, func(items []studentRecord) ([]studentRecord, error) {
		for i := range items {
			if items[i].ID == id.String() {
				return append(items[:i], items[i+1:]...), nil
			}
		}
		return nil, shared.ErrStudentNotFound
	})
}

func (r *StudentRepo) findOnce(ctx context.Context, matchFn func(studentRecord) bool) (*student.Student, error) {
	items, err := r.store.read(ctx)
	if err != nil {
		return nil, err
	}
	for _, rec := range items {
		if matchFn(rec) {
			return toStudent(rec)
		}
	}
	return nil, shared.ErrStudentNotFound
}
