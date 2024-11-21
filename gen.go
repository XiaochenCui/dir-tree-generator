package main

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

type dir struct {
	Path     string   `yaml:"path"`
	Desc     []string `yaml:"desc"`
	Children []dir    `yaml:"children"`
}

// e.g:
// input:
//   - path: "~/.cache/bazel/"
//     desc: "outputRoot"
//     children:
//   - path: "_bazel_<user-name>"
//     desc: "outputUserRoot"
//
// output:
// ~/.cache/bazel/                             <= outputRoot
// └─_bazel_<user-name>/                       <= outputUserRoot
func generate(yamlBytes []byte) (outputBytes []byte, err error) {
	var dirs []dir
	err = yaml.Unmarshal(yamlBytes, &dirs)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	outputBytes = make([]byte, 0, 1024)

	fmt.Printf("root: %+v\n", dirs)

	for _, d := range dirs {
		outputBytes = append(outputBytes, printDir(d, 0)...)
	}

	return outputBytes, nil
}

const dirSpaceWidth = 30

func printDir(d dir, depth int) (outputBytes []byte) {
	var prefix []byte
	var prefixLen int

	if depth > 0 {
		prefix = []byte("└─")

		// takes 6 bytes, but printed as 2 characters
		prefixLen = 2
	}

	// write suffix
	outputBytes = append(outputBytes, prefix...)

	// write path
	outputBytes = append(outputBytes, []byte(d.Path)...)

	// write space
	spaceWidth := dirSpaceWidth - prefixLen - len(d.Path)
	fmt.Printf("spaceWidth: %d, len(suffix): %d, len(d.Path): %d\n", spaceWidth, len(prefix), len(d.Path))
	outputBytes = append(outputBytes, bytes.Repeat([]byte(" "), spaceWidth)...)

	// write desc
	outputBytes = append(outputBytes, []byte("<=")...)
	for i := 0; i < len(d.Desc); i++ {
		outputBytes = append(outputBytes, []byte(d.Desc[i])...)
	}

	outputBytes = append(outputBytes, []byte("\n")...)

	for _, c := range d.Children {
		outputBytes = append(outputBytes, printDir(c, depth+1)...)
	}

	return outputBytes
}
