package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <filename>")
		os.Exit(1)
	}

	filepath := os.Args[1]
	yamlBytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gotOutputBytes, err := Generate(yamlBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(gotOutputBytes))
}
