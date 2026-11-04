package rubric

import "pos/internal/domain/shared"

type Builder struct {
	title  string
	drafts []criterionDraft
}

type criterionDraft struct {
	name     string
	maxScore float64
	weight   float64
}

func NewBuilder(title string) *Builder {
	return &Builder{title: title}
}

func (b *Builder) WithCriterion(name string, maxScore, weight float64) *Builder {
	b.drafts = append(b.drafts, criterionDraft{name: name, maxScore: maxScore, weight: weight})
	return b
}

func (b *Builder) Build() (*Rubric, error) {
	if len(b.drafts) == 0 {
		return nil, shared.ErrEmptyRubric
	}
	r, err := New(b.title)
	if err != nil {
		return nil, err
	}
	for _, d := range b.drafts {
		max, err := NewMaxScore(d.maxScore)
		if err != nil {
			return nil, err
		}
		w, err := NewWeight(d.weight)
		if err != nil {
			return nil, err
		}
		if _, err := r.AddCriterion(d.name, max, w); err != nil {
			return nil, err
		}
	}
	return r, nil
}