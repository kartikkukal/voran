package main

import (
	"fmt"
	"language/module/emit"
	"language/module/gen"
	"language/module/lexer"
	"language/module/node"
	"language/module/parse"
	"language/module/token"
	"os"
)

func polishNotation(node *node.Node) string {
	if node.Token.Kind != token.IntConstant {

		var op string

		switch node.Token.Kind {
		case token.Add:
			op = "+"
		case token.Subtract:
			op = "-"
		case token.Multiply:
			op = "*"
		case token.Divide:
			op = "/"
		case token.Modulus:
			op = "%"
		case token.LessThan:
			op = "<"
		case token.GreaterThan:
			op = ">"
		case token.LessOrEqual:
			op = "<="
		case token.GreaterOrEqual:
			op = ">="
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
	}

	return string(node.Token.Data)
}

func main() {
	data, err := os.ReadFile("tests/exp.v")
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
		t.Debug()
	}

	parserObject := parse.New(tokens)
	n := parserObject.ParseExpression(0)

	p := polishNotation(n)
	fmt.Println(p)

	genISA := gen.NewVoran()
	emit := emit.New(genISA)

	_, err = emit.EmitExpression(n)
	if err != nil {
		fmt.Println(err)
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
