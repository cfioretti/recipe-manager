package domain

type Pans struct {
	Pans             []Pan
	TotalDoughWeight float64
}

type Pan struct {
	Shape       string
	Measures    Measures
	DoughWeight float64
}

type Measures struct {
	Diameter *int
	Edge     *int
	Width    *int
	Length   *int
}

type Strategy func(data map[string]string) Pan
