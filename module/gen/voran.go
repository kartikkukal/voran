package gen

import "fmt"

type Voran struct {
	instructions []string
	registers    []string
}

func NewVoran() *Voran {
	return &Voran{
		instructions: []string{},
		registers:    []string{"q10", "q9", "q8", "q7", "q6", "q5", "q4", "q3", "q2", "q1", "q0"},
	}
}

func (self *Voran) append(format string, a ...any) {
	i := fmt.Sprintf(format, a...)
	self.instructions = append(self.instructions, i)
}

func (self *Voran) AllocateRegister() (string, bool) {

	if len(self.registers) < 1 {
		return "", false
	}

	r := self.registers[0]
	self.registers = self.registers[1:]

	return r, true
}

func (self *Voran) DeallocateRegister(r string) bool {
	self.registers = append(self.registers, r)
	return true
}

func (self *Voran) Load(r string, v int) {
	self.append("ld %v, $%v", r, v)
}

func (self *Voran) Add(r1 string, r2 string) {
	self.append("uadd %v, %v", r1, r2)
}

func (self *Voran) Subtract(r1 string, r2 string) {
	self.append("usub %v, %v", r1, r2)
}

func (self *Voran) Multiply(r1 string, r2 string) {
	self.append("umul %v, %v", r1, r2)
}

func (self *Voran) Divide(r1 string, r2 string) {
	self.append("udiv %v, %v", r1, r2)
}

func (self *Voran) Modulus(r1 string, r2 string) {
	self.append("umod %v, %v", r1, r2)
}

func (self *Voran) Equality(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("seteq %v, $1", r1)
}

func (self *Voran) NotEqual(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("setne %v, $1", r1)
}

func (self *Voran) LessThan(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("setlt %v, $1", r1)
}

func (self *Voran) GreaterThan(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("setgt %v, $1", r1)
}

func (self *Voran) LessOrEqual(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("setle %v, $1", r1)
}

func (self *Voran) GreaterOrEqual(r1 string, r2 string) {
	self.append("ucmp %v, %v", r1, r2)
	self.append("setge %v, $1", r1)
}

func (self Voran) GetInstructions() []string {

	var instructions []string

	preamble := []string{
		".main:",
	}

	postamble := []string{
		"int $3",
		"int $0",
	}

	instructions = append(instructions, preamble...)
	instructions = append(instructions, self.instructions...)
	instructions = append(instructions, postamble...)

	return instructions
}
