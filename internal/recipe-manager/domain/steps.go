package domain

type Steps struct {
	RecipeId int
	Steps    []Step
}

type Step struct {
	Id          int
	StepNumber  int
	Description string
}
