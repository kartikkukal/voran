package gen

import (
	"errors"
	"fmt"
	"slices"
)

var (
	ErrRegisterSpilled = errors.New("register spillover not implemented yet")
)

type Voran struct {
	registers []string
	size      []string
	labels    map[string][]string
	stack     int

	currentFunction string
	currentLabel    string
	labelCounter    int
}

func NewVoran() *Voran {
	return &Voran{
		registers: []string{"12", "11", "10", "9", "8", "7", "6", "5", "4", "3", "2", "1", "0"},
		size:      []string{"b", "w", "d", "q"},

		labels: make(map[string][]string),

		currentFunction: "main",
		currentLabel:    "main",

		stack: 0,
	}
}

func (self *Voran) append(format string, a ...any) {
	self.labels[self.currentLabel] = append(self.labels[self.currentLabel], fmt.Sprintf("\t"+format, a...))
}

func (self *Voran) AllocateRegister(bytes int) (string, bool) {

	r := self.registers[0]
	self.registers = self.registers[1:]

	size := 0

	switch bytes {
	case 1:
		size = 0
	case 2:
		size = 1
	case 4:
		size = 2
	case 8:
		size = 3
	}

	return fmt.Sprintf("%v%v", self.size[size], r), true
}

func (self *Voran) SizedRegister(r string, bytes int) string {
	original := r[1:]

	size := 0

	switch bytes {
	case 1:
		size = 0
	case 2:
		size = 1
	case 4:
		size = 2
	case 8:
		size = 3
	}

	return fmt.Sprintf("%v%v", self.size[size], original)
}

func (self *Voran) DeallocateRegister(r string) bool {

	if len(r) < 2 || len(r) > 3 {
		return false
	}

	size := string(r[0])

	if !slices.Contains(self.size, size) {
		return false
	}

	self.registers = append(self.registers, r[1:])

	return true
}

func (self *Voran) DeallocateAllRegisters() {
	self.registers = []string{"12", "11", "10", "9", "8", "7", "6", "5", "4", "3", "2", "1", "0"}
}

func (self *Voran) CreateLabel() string {
	name := fmt.Sprintf("%v_%v", self.currentFunction, self.labelCounter)
	self.labelCounter++

	self.labels[name] = make([]string, 0)

	return name
}

func (self *Voran) SwitchLabel(name string) {

	if name == "main" {
		self.currentFunction = "main"
	}

	_, ok := self.labels[name]
	if !ok {
		self.labels[name] = make([]string, 0)
	}

	self.currentLabel = name
}

func (self *Voran) Add(r1 string, r2 string) {
	self.append("uadd %v, %v", r1, r2)
	self.DeallocateRegister(r2)
}

func (self *Voran) Subtract(r1 string, r2 string) {
	self.append("usub %v, %v", r1, r2)
	self.DeallocateRegister(r2)
}

func (self *Voran) Multiply(r1 string, r2 string) {
	self.append("umul %v, %v", r1, r2)
	self.DeallocateRegister(r2)
}

func (self *Voran) Divide(r1 string, r2 string) {
	self.append("udiv %v, %v", r1, r2)
	self.DeallocateRegister(r2)
}

func (self *Voran) Modulus(r1 string, r2 string) {
	self.append("umod %v, %v", r1, r2)
	self.DeallocateRegister(r2)
}

func (self *Voran) compareLoad(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("ld %v, 0", r1)
}

func (self *Voran) Equality(r1 string, r2 string) {
	self.compareLoad(r1, r2)
	self.append("seteq %v, 1", r1)
	self.DeallocateRegister(r2)
}

func (self *Voran) NotEqual(r1 string, r2 string) {
	self.compareLoad(r1, r2)
	self.append("setne %v, 1", r1)
	self.DeallocateRegister(r2)
}

func (self *Voran) LessThan(r1 string, r2 string) {
	self.compareLoad(r1, r2)
	self.append("setlt %v, 1", r1)
	self.DeallocateRegister(r2)
}

func (self *Voran) GreaterThan(r1 string, r2 string) {
	self.compareLoad(r1, r2)
	self.append("setgt %v, 1", r1)
	self.DeallocateRegister(r2)
}

func (self *Voran) LessOrEqual(r1 string, r2 string) {
	self.compareLoad(r1, r2)
	self.append("setle %v, 1", r1)
	self.DeallocateRegister(r2)
}

func (self *Voran) GreaterOrEqual(r1 string, r2 string) {
	self.compareLoad(r1, r2)
	self.append("setge %v, 1", r1)
	self.DeallocateRegister(r2)
}

func (self *Voran) Compare(r1 string) {
	self.append("ucmp %v, 1", r1)
}

func (self *Voran) BranchNotEqual(name string) {
	self.append("bne %v", name)
}

func (self *Voran) BranchEqual(name string) {
	self.append("beq %v", name)
}

func (self *Voran) Jump(name string) {
	self.append("jmp %v", name)
}

// Use adequate register sizes
func (self *Voran) LoadLiteral(v, size int) string {
	r, ok := self.AllocateRegister(size)
	if !ok {
		return ""
	}

	self.append("ld %v, 0x%x", r, v)
	return r
}

func (self *Voran) DeclareLocal(size int) int {
	index := self.stack
	self.stack += size

	return index
}

func (self *Voran) LoadLocal(index, size int) string {
	r, ok := self.AllocateRegister(size)
	if !ok {
		return ""
	}
	self.append("ld %v, [q13+%v]", r, index)
	return r
}

func (self *Voran) StoreLocal(r string, index int) {
	self.append("str %v, [q13+%v]", r, index)
}

func (self *Voran) GetAddress(index int) string {
	r, ok := self.AllocateRegister(8)
	if !ok {
		return ""
	}

	self.append("mov %v, q13+%v", r, index)
	return r
}

func (self *Voran) LoadAddress(r string, size int) string {
	r1, ok := self.AllocateRegister(size)
	if !ok {
		return ""
	}

	self.append("ld %v, [%v]", r1, r)
	return r1
}

func (self *Voran) DeclareFunction(name string) {

	self.labelCounter = 0
	self.currentFunction = name

	self.SwitchLabel(name)
	self.append("str q13, [q14+8]")

	self.stack = 16
	self.append("mov q13, q14")
}

func (self *Voran) ReturnFunction() {

	r, ok := self.AllocateRegister(3)
	if !ok {
		fmt.Println("errr")
	}

	self.append("ld %v, [q13+0]", r)
	self.append("ld q13, [q13+8]")

	self.append("jmp %v, 12", r)

	self.stack = 0
}

func (self *Voran) push(r string) {
	self.append("str %v, [q13+%v]", r, self.stack)
	self.stack += 8
}

func (self *Voran) CallFunction(name string) {

	self.append("mov q14, q13")
	self.append("uadd q14, %v", self.stack)
	self.append("str q15, [q14]")

	self.Jump(name)
}

/*
	func (self *Voran) DeclareInt(name string) {
		self.symbols[name] = self.stack
		self.stack += 8
	}

	func (self *Voran) LoadInt(name string) string {
		r, ok := self.AllocateRegister(3)
		if !ok {
			return ""
		}

		location := self.symbols[name]

		self.append("ld %v, [q13+%v]", r, location)
		return r
	}

	func (self *Voran) StoreInt(r, name string) {
		location := self.symbols[name]

		self.append("str %v, [q13+%v]", r, location)
		self.DeallocateRegister(r)
	}
*/
func (self *Voran) Interrupt(i int) {
	self.append("int %v", i)
}

func (self Voran) GetInstructions() []string {

	self.append("int 3")
	self.append("int 0")

	bss := []string{
		"section .bss",
	}
	/*
		for _, v := range self.variables {
			i := fmt.Sprintf("\t%v: res 8", v)
			bss = append(bss, i)
		}
	*/
	var instructions []string

	stack_size := fmt.Sprintf("%%stack: %v", 1024)

	instructions = append(instructions, stack_size)

	instructions = append(instructions, bss...)

	preamble := []string{
		"section .text",
	}

	instructions = append(instructions, preamble...)

	for key, value := range self.labels {
		header := fmt.Sprintf("%v:", key)

		instructions = append(instructions, header)
		instructions = append(instructions, value...)
	}

	return instructions
}
