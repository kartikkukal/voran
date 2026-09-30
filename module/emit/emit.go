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

	sizes = map[ast.Type]int{
		ast.Byte:  1,
		ast.Short: 2,
		ast.Int:   4,
		ast.Long:  8,
	}
)

type Symbol struct {
	dataType ast.Type
	size     int
	address  int
}

type Emit struct {
	locals  map[int]Symbol
	mapping []string

	gen.Gen
}

func New(generator gen.Gen, mapping []string) Emit {

	object := Emit{
		locals: make(map[int]Symbol),

		mapping: mapping,
		Gen:     generator,
	}

	return object
}

func (self *Emit) getName(index int) string {

	if index >= len(self.mapping) {

		// Proper error handling
		return ""
	}

	return self.mapping[index]
}

func (self *Emit) registerSymbol(index int, t ast.Type, address int) {

	self.locals[index] = Symbol{
		dataType: t,
		size:     sizes[t],
		address:  address,
	}
}

func (self *Emit) EmitArguments(master *ast.Node) {

	fmt.Println(master.Debug())

	if master.Left.Kind == ast.Identifier {
		index := self.locals[master.Left.Value]
		_ = self.LoadLocal(index.address, index.size)
	} else {
		self.EmitArguments(master.Left)
	}

	if master.Right.Kind == ast.Identifier {
		index := self.locals[master.Right.Value]
		_ = self.LoadLocal(index.address, index.size)
	} else {
		self.EmitArguments(master.Right)
	}
}

func (self *Emit) EmitExpression(master *ast.Node, register string) (string, error) {

	if master.Kind == ast.Call {
		self.DeallocateAllRegisters()
		self.EmitArguments(master.Left)

		self.CallFunction(self.mapping[master.Value])
		return "", nil
	}

	if master.Kind == ast.Literal {
		// TODO: Add proper handling for larger integer sizes

		r := self.LoadLiteral(master.Value, sizes[master.Type])
		return r, nil
	}

	if master.Kind == ast.Identifier {
		if master.RValue {
			symbol := self.locals[master.Value]
			r := self.LoadLocal(symbol.address, symbol.size)
			return r, nil

		} else {
			// FIXME: Ugly hack to resize register to target size
			register = self.SizedRegister(register, self.locals[master.Value].size)
			self.StoreLocal(register, self.locals[master.Value].address)
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

	switch master.Kind {
	case ast.Equal:

		return right, nil

	case ast.Add:

		self.Add(left, right)
		return left, nil

	case ast.Sub:

		self.Subtract(left, right)
		return left, nil

	case ast.Mul:

		self.Multiply(left, right)
		return left, nil

	case ast.Div:

		self.Divide(left, right)
		return left, nil

	case ast.Mod:

		self.Modulus(left, right)
		return left, nil

	case ast.Equals:

		self.Equality(left, right)
		return left, nil

	case ast.NotEq:

		self.NotEqual(left, right)
		return left, nil

	case ast.Less:

		self.LessThan(left, right)
		return left, nil

	case ast.Greater:

		self.GreaterThan(left, right)
		return left, nil

	case ast.LessEq:

		self.LessOrEqual(left, right)
		return left, nil

	case ast.GreaterEq:

		self.GreaterOrEqual(left, right)
		return left, nil
	}

	return "", ErrUnknownAST
}

func (self *Emit) TraverseParameters(master *ast.Node) error {
	store := func(value int, t ast.Type) error {
		r, ok := self.AllocateRegister(sizes[t])
		if !ok {
			return ErrNoFreeRegisters
		}

		index := self.DeclareLocal(sizes[t])
		self.StoreLocal(r, index)

		self.registerSymbol(value, t, index)
		return nil
	}

	if master.Left.Kind == ast.Declare {
		err := store(master.Left.Value, master.Left.Type)
		if err != nil {
			return err
		}
	} else {
		err := self.TraverseParameters(master.Left)
		if err != nil {
			return err
		}
	}

	if master.Right.Kind == ast.Declare {
		err := store(master.Right.Value, master.Right.Type)
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

	switch master.Kind {
	case ast.Declare:

		size := sizes[master.Type]

		println("Found size: ", size, "\n")

		index := self.DeclareLocal(size)
		self.registerSymbol(master.Value, master.Type, index)

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
		self.locals = make(map[int]Symbol)

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

	case ast.Branch:

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

	/*case ast.Interrupt:

	self.Interrupt(master.Value)
	return nil*/

	default:
		_, err := self.EmitExpression(master, "")
		return err
	}
}
