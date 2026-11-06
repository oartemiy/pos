package rubric

import (
	"pos/internal/domain/shared"
	"slices"
	"strings"
)

type Rubric struct {
	id       shared.ID
	title    string
	criteria []Criterion
}

func New(title string) (*Rubric, error) {
	t := strings.TrimSpace(title)
	if t == "" {
		return nil, shared.ErrEmptyRubricTitle
	}
	return &Rubric{id: shared.NewID(), title: t}, nil
}

func Restore(id shared.ID, title string, criteria []Criterion) (*Rubric, error) {
	r := &Rubric{id: id, title: title}
	r.criteria = make([]Criterion, 0, len(criteria))
	for _, c := range criteria {
		if r.hasName(c.name, c.id) {
			return nil, shared.ErrDuplicateCriterion
		}
		r.criteria = append(r.criteria, c)
	}
	return r, nil
}

func (r *Rubric) ID() shared.ID {
	return r.id
}

func (r *Rubric) Title() string {
	return r.title
}

// NOTE: returns the copy of r.criteria
func (r *Rubric) Criteria() []Criterion {
	return slices.Clone(r.criteria)
}

func (r *Rubric) Criterion(id shared.ID) (Criterion, error) {
	idx := r.indexOf(id)
	if idx < 0 {
		return Criterion{}, shared.ErrCriterionNotFound
	}
	return r.criteria[idx], nil
}

func (r *Rubric) AddCriterion(name string, max MaxScore, w Weight) (Criterion, error) {
	n, err := normalizeCriterionName(name)
	if err != nil {
		return Criterion{}, err
	}
	if r.hasName(n, shared.ID{}) {
		return Criterion{}, shared.ErrDuplicateCriterion
	}
	c := Criterion{id: shared.NewID(), name: n, maxScore: max, weight: w}
	r.criteria = append(r.criteria, c)
	return c, nil
}

func (r *Rubric) RenameCriterion(id shared.ID, name string) error {
	i := r.indexOf(id)
	if i < 0 {
		return shared.ErrCriterionNotFound
	}
	n, err := normalizeCriterionName(name)
	if err != nil {
		return err
	}
	if r.hasName(n, id) {
		return shared.ErrDuplicateCriterion
	}
	r.criteria[i].name = n
	return nil
}

func (r *Rubric) ChangeWeight(id shared.ID, w Weight) error {
	i := r.indexOf(id)
	if i < 0 {
		return shared.ErrCriterionNotFound
	}
	r.criteria[i].weight = w
	return nil
}

func (r *Rubric) RemoveCriterion(id shared.ID) error {
	i := r.indexOf(id)
	if i < 0 {
		return shared.ErrCriterionNotFound
	}
	r.criteria = slices.Delete(r.criteria, i, i+1)
	return nil
}

func (r *Rubric) hasName(name string, id shared.ID) bool {
	return slices.ContainsFunc(r.criteria, func(c Criterion) bool {
		return c.id != id && strings.EqualFold(c.name, name)
	})
}

func (r *Rubric) indexOf(id shared.ID) int {
	return slices.IndexFunc(r.criteria, func(c Criterion) bool { return c.id == id })
}
