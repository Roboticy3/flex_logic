package main

import (
	"C"
	"flex-logic/lcircuit"
	"fmt"
)

const BANNER = `
welcome
			O			\
						 |
			O			/
`

//export TestImport
func TestImport() int {
	return 47
}

func main() {
	fmt.Print(BANNER)

	circuit := lcircuit.CreateCircuit[int, int]()

	fmt.Printf("Successfully created empty circuit at %p", circuit)
}
