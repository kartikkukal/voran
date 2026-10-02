package gen

type Gen interface {
	AllocateRegister(bytes int) (string, bool)
	DeallocateRegister(r string) bool
	SizedRegister(r string, bytes int) string
	DeallocateAllRegisters()

	CreateLabel() string
	SwitchLabel(name string)

	//MainFunction()
	DeclareFunction(name string)
	ReturnFunction()

	CallFunction(name string)

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

	LoadLiteral(v, size int) string

	DeclareLocal(size int) int
	LoadLocal(index, size int) string
	StoreLocal(r string, index int)

	GetAddress(index int) string
	LoadAddress(r string, size int) string

	/*
		DeclareGlobal(name string, size int)
		LoadGlobal(name string)
		StoreGlobal(r, name string)
	*/

	Interrupt(i int)

	GetInstructions() []string
}
