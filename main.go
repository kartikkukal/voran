package main

import (
	"fmt"
	"language/module/emit"
	"language/module/gen"
	"language/module/lexer"
	"language/module/parse"
	"os"
)

/*
func polishNotation(node *ast.Node) string {

	if node.Type == ast.LiteralInt {
		return strconv.Itoa(node.Value)
	}

	if node.Type == ast.Identifier {
		return strconv.Itoa(node.Value)
	}

	var op string

	switch node.Type {
	case ast.Add:
		op = "+"
	case ast.Subtract:
		op = "-"
	case ast.Multiply:
		op = "*"
	case ast.Divide:
		op = "/"
	case ast.Modulus:
		op = "%"
	case ast.Equal:
		op = "="
	default:
		op = "?"
	}

	var left, right string

	if node.Left != nil {
		left = polishNotation(node.Left)
	}
	if node.Right != nil {
		right = polishNotation(node.Right)
	}

	return fmt.Sprintf("(%v %v %v)", op, left, right)

}*/

func main() {
	data, err := os.ReadFile("tests/length.v")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	lexerObject := lexer.New(data)
	tokens, err := lexerObject.Tokenize()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, t := range tokens {
		fmt.Println(t)
	}

	parserObject := parse.New(tokens)
	n, err := parserObject.ParseStatements()
	if err != nil {
		fmt.Println(err)
		return
	}
	/*p := polishNotation(n)
	fmt.Println(p)*/

	genISA := gen.NewVoran()
	emit := emit.New(genISA, parserObject.Symbols)

	fmt.Println(n)

	err = emit.EmitStatement(n)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	f, err := os.Create("a.S")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer f.Close()

	for _, l := range genISA.GetInstructions() {
		fmt.Println(l)

		_, err := f.WriteString(l + "\n")

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
