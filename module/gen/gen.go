package gen

type Gen interface {
	AllocateRegister(size int) (string, bool)
	DeallocateRegister(r string) bool
	DeallocateAllRegisters()

	CreateLabel() string
	SwitchLabel(name string)

	LoadValue(v int) string
	LoadInt(name string) string
	Add(r1, r2 string)
	Subtract(r1, r2 string)
	Multiply(r1, r2 string)
	Divide(r1, r2 string)
	Modulus(r1, r2 string)

	Equality(r1, r2 string)
	NotEqual(r1, r2 string)

	LessThan(r1, r2 string)
	GreaterThan(r1, r2 string)
	LessOrEqual(r1, r2 string)
	GreaterOrEqual(r1, r2 string)

	Compare(r1 string)

	BranchNotEqual(name string)
	BranchEqual(name string)

	Jump(name string)

	DeclareInt(name string)
	StoreInt(r, name string)

	Interrupt(i int)

	GetInstructions() []string
}
