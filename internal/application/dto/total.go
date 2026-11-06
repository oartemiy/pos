package dto

type TotalRow struct {
	StudentID   string
	StudentName string
	Email       string
	Value       float64
	Graded      int
	Criteria    int
	Strategy    string
}

type ExportResult struct {
	Location string
	Rows     int
}

type StudentReport struct {
	Student    StudentView
	Evaluation []EvaluationView
	Total      TotalRow
}
