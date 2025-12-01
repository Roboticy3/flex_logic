package main

import (
	"flex-logic/lcircuit"
	"fmt"
)

const BANNER = `
welcome
			O			\
						 |
			O			/
`

func TestImport() int {
	return 42
}

func main() {
	fmt.Print(BANNER)

	circuit := lcircuit.CreateCircuit[int, int]()

	fmt.Printf("Successfully created empty circuit at %p", circuit)
}
