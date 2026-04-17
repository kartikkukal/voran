package emit

import (
	"errors"
	"language/module/gen"
	n "language/module/node"
	t "language/module/token"
)

var (
	ErrNoValue         = errors.New("terminal node has no value")
	ErrNoFreeRegisters = errors.New("no free registers")
)

type Emit struct {
	gen.Gen
}

func New(generator gen.Gen) Emit {
	return Emit{
		Gen: generator,
	}
}

func (self *Emit) EmitExpression(master *n.Node) (string, error) {

	if master.IsTerminal() {
		r, ok := self.AllocateRegister()
		if !ok {
			return "", ErrNoFreeRegisters
		}

		n, ok := master.Value()
		if !ok {
			return "", ErrNoValue
		}

		self.Load(r, n)
		return r, nil
	}

	left, err := self.EmitExpression(master.Left)
	if err != nil {
		return "", err
	}

	right, err := self.EmitExpression(master.Right)
	if err != nil {
		return "", err
	}

	switch master.Token.Kind {
	case t.Add:

		self.Add(left, right)
		self.DeallocateRegister(right)

		return left, nil

	case t.Subtract:

		self.Subtract(left, right)
		self.DeallocateRegister(right)

		return left, nil

	case t.Multiply:

		self.Multiply(left, right)
		self.DeallocateRegister(right)

		return left, nil

	case t.Divide:

		self.Divide(left, right)
		self.DeallocateRegister(right)

		return left, nil

	case t.Modulus:

		self.Modulus(left, right)
		self.DeallocateRegister(right)

		return left, nil

	case t.Equality:

		self.Equality(left, right)
		self.DeallocateRegister(right)

		return left, nil

	case t.NotEqual:

		self.NotEqual(left, right)
		self.DeallocateRegister(right)

		return left, nil

	case t.LessThan:

		self.LessThan(left, right)
		self.DeallocateRegister(right)

		return left, nil

	case t.GreaterThan:

		self.GreaterThan(left, right)
		self.DeallocateRegister(right)

		return left, nil

	case t.LessOrEqual:

		self.LessOrEqual(left, right)
		self.DeallocateRegister(right)

		return left, nil

	case t.GreaterOrEqual:

		self.GreaterOrEqual(left, right)
		self.DeallocateRegister(right)

		return left, nil
	}

	panic("unknown operator")
}
