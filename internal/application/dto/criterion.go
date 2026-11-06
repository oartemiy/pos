package dto

type AddCriterionInput struct {
	Name     string
	MaxScore float64
	Weight   float64
}

type CriterionView struct {
	ID        string
	Name      string
	MaxWeight float64
	Weight    float64
}
