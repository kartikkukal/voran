package main

import (
	"flag"
	"fmt"
	"language/module/emit"
	"language/module/gen"
	"language/module/lexer"
	"language/module/parse"
	"os"
	"path/filepath"
	"strings"
)

var (
	sourcePath string
	outputPath string

	tokensDebug   bool
	parserDebug   bool
	assemblyDebug bool
)

func init() {
	flag.StringVar(&sourcePath, "src", "", "path of the source file")
	flag.StringVar(&outputPath, "output", "", "path of the output file")

	flag.BoolVar(&tokensDebug, "debug-tokens", false, "print tokens")
	flag.BoolVar(&parserDebug, "debug-parser", false, "print parser tree")
	flag.BoolVar(&assemblyDebug, "debug-assembly", false, "print assembly output")
}

func main() {

	flag.Parse()

	defaultPath := strings.TrimSuffix(sourcePath, filepath.Ext(sourcePath))

	if outputPath == "" {
		outputPath = defaultPath
	}

	data, err := os.ReadFile(sourcePath)
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

	if tokensDebug {
		for _, t := range tokens {
			fmt.Println(t.Debug())
		}
	}

	parserObject := parse.New(tokens)

	n, err := parserObject.ParseStatements()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if parserDebug {
		fmt.Println(n.Debug())
	}

	genISA := gen.NewVoran()
	emit := emit.New(genISA, parserObject.Symbols)

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

		if assemblyDebug {
			fmt.Println(l)
		}

		_, err := f.WriteString(l + "\n")

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
