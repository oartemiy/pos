package jsonfile

import (
	"fmt"
	"pos/internal/domain/evaluation"
	"pos/internal/domain/rubric"
	"pos/internal/domain/shared"
	"pos/internal/domain/student"
	"time"
)

type studentRecord struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic,omitempty"`
	Email      string `json:"email"`
}

func toStudentRecord(s *student.Student) studentRecord {
	n := s.Name()
	return studentRecord{
		ID:         s.ID().String(),
		LastName:   n.Last(),
		FirstName:  n.First(),
		Patronymic: n.Patronymic(),
		Email:      s.Email().String(),
	}
}

func toStudent(r studentRecord) (*student.Student, error) {
	id, err := shared.ParseID(r.ID)
	if err != nil {
		return nil, fmt.Errorf("jsonfile: student %q: %w", r.ID, err)
	}
	name, err := shared.NewName(r.FirstName, r.LastName, r.Patronymic)
	if err != nil {
		return nil, fmt.Errorf("jsonfile: student %q: %w", r.ID, err)
	}
	email, err := shared.NewEmail(r.Email)
	if err != nil {
		return nil, fmt.Errorf("jsonfile: student %q: %w", r.ID, err)
	}
	return student.Restore(id, name, email), nil
}

type criterionRecord struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	MaxScore float64 `json:"max_score"`
	Weight   float64 `json:"weight"`
}

type rubricRecord struct {
	ID       string            `json:"id"`
	Title    string            `json:"title"`
	Criteria []criterionRecord `json:"criteria"`
}

func toRubricRecord(r *rubric.Rubric) rubricRecord {
	criteria := r.Criteria()
	recs := make([]criterionRecord, 0, len(criteria))
	for _, c := range criteria {
		recs = append(recs, criterionRecord{
			ID:       c.ID().String(),
			Name:     c.Name(),
			MaxScore: c.MaxScore().Value(),
			Weight:   c.Weight().Value(),
		})
	}
	return rubricRecord{ID: r.ID().String(), Title: r.Title(), Criteria: recs}
}

func toRubric(rec rubricRecord) (*rubric.Rubric, error) {
	id, err := shared.ParseID(rec.ID)
	if err != nil {
		return nil, fmt.Errorf("jsonfile: rubric %q: %w", rec.ID, err)
	}

	criteria := make([]rubric.Criterion, 0, len(rec.Criteria))
	for _, cr := range rec.Criteria {
		c, err := toCriterion(cr)
		if err != nil {
			return nil, fmt.Errorf("jsonfile: rubric %q: %w", rec.ID, err)
		}
		criteria = append(criteria, c)
	}

	r, err := rubric.Restore(id, rec.Title, criteria)
	if err != nil {
		return nil, fmt.Errorf("jsonfile: rubric %q: %w", rec.ID, err)
	}
	return r, nil
}

func toCriterion(cr criterionRecord) (rubric.Criterion, error) {
	id, err := shared.ParseID(cr.ID)
	if err != nil {
		return rubric.Criterion{}, fmt.Errorf("criterion %q: %w", cr.ID, err)
	}
	max, err := rubric.NewMaxScore(cr.MaxScore)
	if err != nil {
		return rubric.Criterion{}, fmt.Errorf("criterion %q: %w", cr.ID, err)
	}
	w, err := rubric.NewWeight(cr.Weight)
	if err != nil {
		return rubric.Criterion{}, fmt.Errorf("criterion %q: %w", cr.ID, err)
	}
	return rubric.RestoreCriterion(id, cr.Name, max, w)
}

type evaluationRecord struct {
	ID          string    `json:"id"`
	StudentID   string    `json:"student_id"`
	CriterionID string    `json:"criterion_id"`
	Score       float64   `json:"score"`
	Comment     string    `json:"comment,omitempty"`
	GradedAt    time.Time `json:"graded_at"`
}

func toEvaluationRecord(e *evaluation.Evaluation) evaluationRecord {
	return evaluationRecord{
		ID:          e.ID().String(),
		StudentID:   e.StudentID().String(),
		CriterionID: e.CriterionID().String(),
		Score:       e.Score().Value(),
		Comment:     e.Comment().String(),
		GradedAt:    e.GradedAt(),
	}
}

func toEvaluation(r evaluationRecord) (*evaluation.Evaluation, error) {
	wrap := func(err error) error {
		return fmt.Errorf("jsonfile: evaluation %q: %w", r.ID, err)
	}

	id, err := shared.ParseID(r.ID)
	if err != nil {
		return nil, wrap(err)
	}
	studentID, err := shared.ParseID(r.StudentID)
	if err != nil {
		return nil, wrap(err)
	}
	criterionID, err := shared.ParseID(r.CriterionID)
	if err != nil {
		return nil, wrap(err)
	}
	score, err := evaluation.NewScore(r.Score)
	if err != nil {
		return nil, wrap(err)
	}
	comment, err := evaluation.NewComment(r.Comment)
	if err != nil {
		return nil, wrap(err)
	}
	return evaluation.Restore(id, studentID, criterionID, score, comment, r.GradedAt), nil
}
