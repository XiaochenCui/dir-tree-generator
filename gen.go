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

	ancestorLines := make([]bool, dirWidthLimit)
	for i, d := range dirs {
		isLastChild := i == len(dirs)-1
		outputBytes = append(outputBytes, printDir(d, 0, isLastChild, 0, ancestorLines)...)
	}

	return outputBytes, nil
}

const dirWidthLimit = 44
const totalWidthLimit = 80

// # Arguments
// - ancestorLines: Indicates whether a line should be drawn for the ancestor at each position.
// - start: Start position of the directory path.
func printDir(d dir, depth int, isLastChild bool, start int, ancestorLines []bool) (outputBytes []byte) {
	if d.Path == "k8-fastbuild/" {
		fmt.Printf("d: %+v\n", d)
	}

	childStart := getChildStart(d.Path, start)

	selfPrefixLen := 0

	// write parent prefix
	parentPrefix, _ := getParentPrefix(ancestorLines)
	outputBytes = append(outputBytes, substr(string(parentPrefix), 0, start)...)

	// write self prefix
	outputBytes = append(outputBytes, getSelfPrefix(isLastChild, start)...)

	// write path
	outputBytes = append(outputBytes, []byte(d.Path)...)

	// write space
	spaceWidth := dirWidthLimit - start - selfPrefixLen - len(d.Path)
	outputBytes = append(outputBytes, bytes.Repeat([]byte(" "), spaceWidth)...)

	// write desc
	outputBytes = append(outputBytes, []byte("<= ")...)
	brokenDesc := make([]string, 0, 2)
	for _, desc := range d.Desc {
		brokenDesc = append(brokenDesc, broke(desc, totalWidthLimit-dirWidthLimit)...)
	}

	for i, desc := range brokenDesc {
		if i == 0 {
			outputBytes = append(outputBytes, []byte(d.Desc[i])...)
		} else {
			// write parent prefix
			parentPrefix, _ := getParentPrefix(ancestorLines)
			outputBytes = append(outputBytes, substr(string(parentPrefix), 0, childStart)...)

			outputBytes = append(outputBytes, []byte("│")...)

			// write spaces
			spaceWidth := dirWidthLimit - childStart + 3 + 1 // +3 for "<= ", -1 for "│"
			outputBytes = append(outputBytes, bytes.Repeat([]byte(" "), spaceWidth)...)

			// write desc
			outputBytes = append(outputBytes, []byte(desc)...)
		}
		outputBytes = append(outputBytes, []byte("\n")...)
	}

	parentHasFollowingSibling := !isLastChild
	for i, c := range d.Children {
		cIsFinalChild := i == len(d.Children)-1
		if parentHasFollowingSibling {
			ancestorLines[start] = true
		} else {
			ancestorLines[start] = false
		}
		outputBytes = append(outputBytes, printDir(c, depth+1, cIsFinalChild, childStart, ancestorLines)...)
	}

	return outputBytes
}

func broke(desc string, limit int) (broken []string) {
	if len(desc) <= limit {
		return []string{desc}
	}
	broken = make([]string, 0, 2)
	for len(desc) > limit {
		// find the last space before the limit
		spaceIdx := limit
		for i := limit; i > 0; i-- {
			if desc[i] == ' ' {
				spaceIdx = i
				break
			}
		}
		broken = append(broken, desc[:spaceIdx])
		desc = desc[spaceIdx+1:]
	}
	broken = append(broken, desc)
	return broken
}

func getParentPrefix(ancestorLines []bool) ([]byte, int) {
	var parentPrefix []byte
	parentPrefixLen := 0
	for i := 0; i < len(ancestorLines); i++ {
		if ancestorLines[i] {
			parentPrefix = append(parentPrefix, []byte("│")...)
			parentPrefixLen = i
		} else {
			parentPrefix = append(parentPrefix, []byte(" ")...)
		}
	}
	return parentPrefix, parentPrefixLen
}

// get substring by unicode characters
func substr(s string, start, end int) string {
	return string([]rune(s)[start:end])
}

func getChildStart(path string, start int) (childStart int) {
	if start == 0 {
		childStart = start
	} else {
		childStart = start + 2
	}
	for i := len(path) - 2; i >= 0; i-- {
		if path[i] == '/' {
			return childStart + i + 1
		}
	}
	return childStart
}

func getSelfPrefix(isLastChild bool, start int) []byte {
	if start > 0 {
		if isLastChild {
			return []byte("└─")
		}
		return []byte("├─")
	}
	return []byte("")
}
