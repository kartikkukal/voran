package gen

import (
	"fmt"
	"slices"
)

type Voran struct {
	registers []string
	size      []string
	variables []string
	labels    map[string][]string

	currentLabel string
	labelCounter int
}

func NewVoran() *Voran {
	return &Voran{
		registers:    []string{"12", "11", "10", "9", "8", "7", "6", "5", "4", "3", "2", "1", "0"},
		size:         []string{"b", "w", "d", "q"},
		labels:       make(map[string][]string),
		currentLabel: "main",
		labelCounter: 0,
	}
}

func (self *Voran) append(format string, a ...any) {
	self.labels[self.currentLabel] = append(self.labels[self.currentLabel], fmt.Sprintf("\t"+format, a...))
}

func (self *Voran) AllocateRegister(size int) (string, bool) {

	if size < 0 || size > 3 {
		return "", false
	}

	if len(self.registers) < 1 {
		return "", false
	}

	r := self.registers[0]
	self.registers = self.registers[1:]

	return fmt.Sprintf("%v%v", self.size[size], r), true
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
	name := fmt.Sprintf("L%v", self.labelCounter)
	self.labelCounter++

	self.labels[name] = make([]string, 0)

	return name
}

func (self *Voran) SwitchLabel(name string) {
	self.currentLabel = name
}

func (self *Voran) LoadValue(v int) string {
	r, ok := self.AllocateRegister(3)
	if !ok {
		return ""
	}

	self.append("ld %v, 0x%x", r, v)
	return r
}

func (self *Voran) LoadInt(name string) string {
	r, ok := self.AllocateRegister(3)
	if !ok {
		return ""
	}

	self.append("ld %v, [%v]", r, name)
	return r
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

func (self *Voran) Equality(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("ld %v, 0", r1)
	self.append("seteq %v, 1", r1)
	self.DeallocateRegister(r2)
}

func (self *Voran) NotEqual(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("ld %v, 0", r1)
	self.append("setne %v, 1", r1)
	self.DeallocateRegister(r2)
}

func (self *Voran) LessThan(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("ld %v, 0", r1)
	self.append("setlt %v, 1", r1)
	self.DeallocateRegister(r2)
}

func (self *Voran) GreaterThan(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("ld %v, 0", r1)
	self.append("setgt %v, 1", r1)
	self.DeallocateRegister(r2)
}

func (self *Voran) LessOrEqual(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("ld %v, 0", r1)
	self.append("setle %v, 1", r1)
	self.DeallocateRegister(r2)
}

func (self *Voran) GreaterOrEqual(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("ld %v, 0", r1)
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

func (self *Voran) DeclareInt(name string) {
	self.variables = append(self.variables, name)
}

func (self *Voran) StoreInt(r, name string) {
	self.append("str %v, [%v]", r, name)
	self.DeallocateRegister(r)
}

func (self *Voran) Interrupt(i int) {
	self.append("int %v", i)
}

func (self Voran) GetInstructions() []string {

	self.append("int 5")
	self.append("int 0")

	bss := []string{
		"section .bss",
	}

	for _, v := range self.variables {
		i := fmt.Sprintf("\t%v: res 8", v)
		bss = append(bss, i)
	}

	var instructions []string
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
