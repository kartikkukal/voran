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

	if master.Type == ast.LiteralInt {
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
	case ast.Equal:

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

	case ast.IsEqual:

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

	switch master.Type {
	case ast.DeclareInt:

		self.DeclareInt(self.symbols[master.Value])
		return nil

	case ast.IfElse:

		end := self.CreateLabel()

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

		self.Jump(end)
		self.SwitchLabel(end)

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

	default:
		_, err := self.EmitExpression(master, "")
		self.DeallocateAllRegisters()
		return err
	}

	return nil
}
