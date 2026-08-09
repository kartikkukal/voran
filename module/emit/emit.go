package emit

import (
	"errors"
	"fmt"
	"language/module/ast"
	"language/module/gen"
)

var (
	ErrNoValue         = errors.New("terminal node has no value")
	ErrNoFreeRegisters = errors.New("no free registers")
	ErrUnknownAST      = errors.New("unknown AST node encountered")
)

type Symbol struct {
	stack int
}

type Emit struct {
	locals  map[int]int
	globals map[int]string

	mapping []string

	gen.Gen
}

func New(generator gen.Gen, mapping []string) Emit {

	object := Emit{
		locals:  make(map[int]int),
		globals: make(map[int]string, 0),

		mapping: mapping,
		Gen:     generator,
	}

	return object
}

func (self *Emit) EmitArguments(master *ast.Node) {

	fmt.Println(master.Debug())

	if master.Left.Type == ast.Identifier {
		index := self.locals[master.Left.Value]
		_ = self.LoadLocal(index, 8)
	} else {
		self.EmitArguments(master.Left)
	}

	if master.Right.Type == ast.Identifier {
		index := self.locals[master.Right.Value]
		_ = self.LoadLocal(index, 8)
	} else {
		self.EmitArguments(master.Right)
	}
}

func (self *Emit) EmitExpression(master *ast.Node, register string) (string, error) {

	if master.Type == ast.Call {
		self.DeallocateAllRegisters()
		self.EmitArguments(master.Left)

		self.CallFunction(self.mapping[master.Value])
		return "", nil
	}

	if master.Type == ast.Integer {
		r := self.LoadLiteral(master.Value)
		return r, nil
	}

	if master.Type == ast.Identifier {
		if master.RValue {
			r := self.LoadLocal(self.locals[master.Value], 3)
			return r, nil

		} else {
			self.StoreLocal(register, self.locals[master.Value])
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

/*
func (self *Emit)
*/
func (self *Emit) TraverseParameters(master *ast.Node) error {
	store := func(value int) error {
		r, ok := self.AllocateRegister(3)
		if !ok {
			return ErrNoFreeRegisters
		}

		index := self.DeclareLocal(8)
		self.StoreLocal(r, index)

		self.locals[value] = index
		return nil
	}

	if master.Left.Type == ast.Declare {
		err := store(master.Left.Value)
		if err != nil {
			return err
		}
	} else {
		err := self.TraverseParameters(master.Left)
		if err != nil {
			return err
		}
	}

	if master.Right.Type == ast.Declare {
		err := store(master.Right.Value)
		if err != nil {
			return err
		}
	} else {
		err := self.TraverseParameters(master.Right)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Emit) EmitStatement(master *ast.Node) error {

	self.DeallocateAllRegisters()

	switch master.Type {
	case ast.Declare:

		index := self.DeclareLocal(8)
		self.locals[master.Value] = index

		if master.Left != nil {
			self.EmitExpression(master.Left, "")
		}

		return nil

	case ast.Function:

		if self.mapping[master.Value] == "main" {
			self.SwitchLabel("main")
			return self.EmitStatement(master.Right)
		}

		self.DeclareFunction(self.mapping[master.Value])
		self.locals = make(map[int]int)

		self.DeallocateAllRegisters()
		err := self.TraverseParameters(master.Left)
		if err != nil {
			return err
		}

		err = self.EmitStatement(master.Right)
		if err != nil {
			return err
		}

		self.ReturnFunction()

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
