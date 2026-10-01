package parse

import (
	"language/module/ast"
	"language/module/token"
)

var (
	intWeight = map[ast.Type]int{
		ast.Byte:  0,
		ast.Short: 1,
		ast.Int:   2,
		ast.Long:  3,

		ast.UByte:  4,
		ast.UShort: 5,
		ast.UInt:   6,
		ast.ULong:  7,
	}

	floatWeight = map[ast.Type]int{
		ast.Float:  0,
		ast.Double: 1,
	}

	broadType = map[ast.Type]int{
		ast.Byte:  0,
		ast.Short: 0,
		ast.Int:   0,
		ast.Long:  0,

		ast.UByte:  1,
		ast.UShort: 1,
		ast.UInt:   1,
		ast.ULong:  1,

		ast.Float:  2,
		ast.Double: 2,

		ast.String: 3,

		ast.None: 4,
	}
)

func (self *Parser) ParseType(t string) (ast.Type, error) {
	typeMap := map[string]ast.Type{
		token.Byte:  ast.Byte,
		token.Short: ast.Short,
		token.Int:   ast.Int,
		token.Long:  ast.Long,

		token.UByte:  ast.UByte,
		token.UShort: ast.UShort,
		token.UInt:   ast.UInt,
		token.ULong:  ast.ULong,

		token.Float:  ast.Float,
		token.Double: ast.Double,
		token.String: ast.String,
	}

	astType, ok := typeMap[t]
	if !ok {
		return ast.None, ErrExpectedType
	}

	return astType, nil
}

func (self *Parser) CoerceTypeCast(master *ast.Node) (ast.Type, error) {

	lType := ast.None
	rType := ast.None

	if master.Terminal() {
		return master.Type, nil
	}

	var err error

	if master.Left.Terminal() {
		lType = master.Left.Type
	} else {
		lType, err = self.CoerceTypeCast(master.Left)
		if err != nil {
			return ast.None, err
		}
	}

	if master.Right.Terminal() {
		rType = master.Right.Type
	} else {
		rType, err = self.CoerceTypeCast(master.Right)
		if err != nil {
			return ast.None, err
		}
	}

	if broadType[lType] != broadType[rType] {
		// FIXME: Add proper handling for this error
		panic("cannot handle broad type conversion implicitly")
	}

	if broadType[lType] == 0 || broadType[lType] == 1 {
		var cast *ast.Node

		if intWeight[lType] == intWeight[rType] {
			return lType, nil
		}

		if intWeight[lType]-intWeight[rType] > 0 {
			cast = ast.NewWithValue(ast.Implicit, int(rType), lType)
			cast.Left = master.Right
			master.Right = cast
		}

		if intWeight[rType]-intWeight[lType] > 0 {
			cast = ast.NewWithValue(ast.Implicit, int(lType), rType)
			cast.Left = master.Left
			master.Left = cast
		}

		master.Type = cast.Type

		return cast.Type, nil
	}

	// TODO: Implement implicit casting for other types
	panic("implicit resizing for non-integer types not implelemented")
}
