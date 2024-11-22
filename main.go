package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: dir_tree_generator <yaml_file>")
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
