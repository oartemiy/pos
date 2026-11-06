package usecase

import (
	"cmp"
	"context"
	"errors"
	"pos/internal/application/dto"
	"pos/internal/domain/evaluation"
	"pos/internal/domain/shared"
	"pos/internal/domain/student"
	"slices"
)

type Students struct {
	students   student.Repository
	evaluation evaluation.Repository
}

func NewStudents (s student.Repository, e evaluation.Repository) *Students {
	return &Students{s, e}
}

func (uc *Students) Add(ctx context.Context, in dto.AddStudentInput) (dto.StudentView, error) {
	name, err := shared.NewName(in.FirstName, in.LastName, in.Patronymic)
	if err != nil {
		return dto.StudentView{}, err
	}
	email, err := shared.NewEmail(in.Email)
	if err != nil {
		return dto.StudentView{}, err
	}

	_, err = uc.students.FindByEmail(ctx, email)

	switch {
	case err == nil:
		return dto.StudentView{}, shared.ErrEmailTaken
	case !errors.Is(err, shared.ErrNotFound):
		return dto.StudentView{}, err
	}
	s := student.New(name, email)
	if err = uc.students.Save(ctx, s); err != nil {
		return dto.StudentView{}, err
	}
	return toStudentView(s), nil
}

func (uc *Students) Get(ctx context.Context, rawID string) (dto.StudentView, error) {
	id, err := shared.ParseID(rawID)
	if err != nil {
		return dto.StudentView{}, err
	}
	s, err := uc.students.FindByID(ctx, id)
	if err != nil {
		return dto.StudentView{}, notFoundAs(err, shared.ErrStudentNotFound)
	}
	return toStudentView(s), nil
}

func (uc *Students) List(ctx context.Context) ([]dto.StudentView, error) {
	list, err := uc.students.List(ctx)
	if err != nil {
		return nil, err
	}
	sortStudents(list)

	views := make([]dto.StudentView, 0, len(list))
	for _, s := range list {
		views = append(views, toStudentView(s))
	}
	return views, nil
}

// Remove удаляет студента, только если у него нет оценок: иначе они
// остались бы ссылаться на несуществующего студента.
func (uc *Students) Remove(ctx context.Context, rawID string) error {
	id, err := shared.ParseID(rawID)
	if err != nil {
		return err
	}
	evals, err := uc.evaluation.FindByStudent(ctx, id)
	if err != nil {
		return err
	}
	if len(evals) > 0 {
		return shared.ErrStudentHasEvaluations
	}
	return notFoundAs(uc.students.Delete(ctx, id), shared.ErrStudentNotFound)
}

func toStudentView(s *student.Student) dto.StudentView {
	return dto.StudentView{ID: s.ID().String(), FullName: s.Name().String(), ShortName: s.Name().Short(), Email: s.Email().String()}
}

func sortStudents(list []*student.Student) {
	slices.SortFunc(list, func(a, b *student.Student) int {
		return cmp.Compare(a.Name().String(), b.Name().String())
	})
}
