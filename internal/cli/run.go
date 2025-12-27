package cli

import (
	"fmt"
	"os"
	"vsharp/internal/lexer"
)

func run(file string) {
	source, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	tokens, err := lexer.Tokenize(string(source), file)
	if err != nil {
		panic(err)
	}
	for _, ex := range tokens {
		fmt.Println(ex.Type.String())
	}
}
