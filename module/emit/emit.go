package emit

import (
	"errors"
	"language/module/ast"
	"language/module/gen"
)

var (
	ErrNoValue         = errors.New("terminal node has no value")
	ErrNoFreeRegisters = errors.New("no free registers")
	ErrUnknownAST      = errors.New("unknown AST node encountered")
)

type Emit struct {
	symbols []string
	gen.Gen
}

func New(generator gen.Gen, symbols []string) Emit {
	return Emit{
		symbols: symbols,
		Gen:     generator,
	}
}

func (self *Emit) EmitExpression(master *ast.Node, register string) (string, error) {

	if master.Type == ast.Integer {
		r := self.LoadValue(master.Value)
		return r, nil
	}

	if master.Type == ast.Identifier {
		if master.RValue {
			r := self.LoadInt(self.symbols[master.Value])
			return r, nil
		} else {
			self.StoreInt(register, self.symbols[master.Value])
			return register, nil
		}
	}

	left, err := self.EmitExpression(master.Left, "")
	if err != nil {
		return "", err
	}

	right, err := self.EmitExpression(master.Right, left)
	if err != nil {
		return "", err
	}

	switch master.Type {
	case ast.Equals:

		return right, nil

	case ast.Add:

		self.Add(left, right)
		return left, nil

	case ast.Subtract:

		self.Subtract(left, right)
		return left, nil

	case ast.Multiply:

		self.Multiply(left, right)
		return left, nil

	case ast.Divide:

		self.Divide(left, right)
		return left, nil

	case ast.Modulus:

		self.Modulus(left, right)
		return left, nil

	case ast.Equal:

		self.Equality(left, right)
		return left, nil

	case ast.NotEqual:

		self.NotEqual(left, right)
		return left, nil

	case ast.LessThan:

		self.LessThan(left, right)
		return left, nil

	case ast.GreaterThan:

		self.GreaterThan(left, right)
		return left, nil

	case ast.LessOrEqual:

		self.LessOrEqual(left, right)
		return left, nil

	case ast.GreaterOrEqual:

		self.GreaterOrEqual(left, right)
		return left, nil
	}

	return "", ErrUnknownAST
}

func (self *Emit) EmitStatement(master *ast.Node) error {

	self.DeallocateAllRegisters()

	switch master.Type {
	case ast.Declare:

		self.DeclareInt(self.symbols[master.Value])
		return nil

	case ast.If:

		label_1 := self.CreateLabel()

		r, err := self.EmitExpression(master.Left, "")
		if err != nil {
			return err
		}

		self.Compare(r)
		self.BranchNotEqual(label_1)

		err = self.EmitStatement(master.Right.Left)
		if err != nil {
			return err
		}

		if master.Right.Right == nil {
			self.Jump(label_1)
			self.SwitchLabel(label_1)

			return nil
		}

		label_2 := self.CreateLabel()

		self.Jump(label_2)
		self.SwitchLabel(label_1)

		err = self.EmitStatement(master.Right.Right)
		if err != nil {
			return err
		}

		self.Jump(label_2)
		self.SwitchLabel(label_2)

		return nil

	case ast.While:

		start := self.CreateLabel()
		end := self.CreateLabel()

		self.Jump(start)
		self.SwitchLabel(start)

		r, err := self.EmitExpression(master.Left, "")
		if err != nil {
			return err
		}

		self.Compare(r)
		self.BranchNotEqual(end)

		err = self.EmitStatement(master.Right)
		if err != nil {
			return err
		}

		self.Jump(start)
		self.SwitchLabel(end)

		return nil

	case ast.Glue:

		if master.Left != nil {
			err := self.EmitStatement(master.Left)
			if err != nil {
				return err
			}
		}

		if master.Right != nil {
			err := self.EmitStatement(master.Right)
			if err != nil {
				return err
			}
		}

		return nil

	case ast.Interrupt:

		self.Interrupt(master.Value)
		return nil

	default:
		_, err := self.EmitExpression(master, "")
		return err
	}
}
