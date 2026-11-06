package dto

type AddStudentInput struct {
	LastName   string
	FirstName  string
	Patronymic string
	Email      string
}

type StudentView struct {
	ID        string
	FullName  string
	ShortName string
	Email     string
}
