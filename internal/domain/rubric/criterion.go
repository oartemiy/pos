package rubric

import (
	"strings"
	"unicode/utf8"

	"pos/internal/domain/shared"
)

const maxCriterionNameLen = 100

type Criterion struct {
	id       shared.ID
	name     string
	maxScore MaxScore
	weight   Weight
}

func RestoreCriterion(id shared.ID, name string, max MaxScore, w Weight) (Criterion, error) {
	n, err := normalizeCriterionName(name)
	if err != nil {
		return Criterion{}, err
	}
	return Criterion{id: id, name: n, maxScore: max, weight: w}, nil
}

func (c Criterion) ID() shared.ID      { return c.id }
func (c Criterion) Name() string       { return c.name }
func (c Criterion) MaxScore() MaxScore { return c.maxScore }
func (c Criterion) Weight() Weight     { return c.weight }

func normalizeCriterionName(raw string) (string, error) {
	n := strings.Join(strings.Fields(raw), " ")
	if n == "" {
		return "", shared.ErrEmptyCriterionName
	}
	if utf8.RuneCountInString(n) > maxCriterionNameLen {
		return "", shared.ErrCriterionNameTooLong
	}
	return n, nil
}
