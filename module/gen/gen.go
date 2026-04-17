package gen

type Gen interface {
	AllocateRegister() (string, bool)
	DeallocateRegister(r string) bool

	Load(r string, v int)
	Add(r1 string, r2 string)
	Subtract(r1 string, r2 string)
	Multiply(r1 string, r2 string)
	Divide(r1 string, r2 string)
	Modulus(r1 string, r2 string)
	Equality(r1 string, r2 string)
	NotEqual(r1 string, r2 string)
	LessThan(r1 string, r2 string)
	GreaterThan(r1 string, r2 string)
	LessOrEqual(r1 string, r2 string)
	GreaterOrEqual(r1 string, r2 string)

	GetInstructions() []string
}
