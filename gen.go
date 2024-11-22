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

// Generate generates a tree structure from a YAML file.
//
// e.g:
// input:
//
//   - path: "~/.cache/bazel/"
//     desc: "outputRoot"
//     children:
//   - path: "_bazel_<user-name>"
//     desc: "outputUserRoot"
//
// output:
//
//	~/.cache/bazel/                             <= outputRoot
//	         └─_bazel_<user-name>/              <= outputUserRoot
func Generate(yamlBytes []byte) (outputBytes []byte, err error) {
	var dirs []dir
	err = yaml.Unmarshal(yamlBytes, &dirs)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	outputBytes = make([]byte, 0, 1024)

	ancestorLines := make([]bool, dirWidthLimit)
	for i, d := range dirs {
		isLastChild := i == len(dirs)-1
		outputBytes = append(outputBytes, printDir(d, isLastChild, 0, ancestorLines)...)
	}

	return outputBytes, nil
}

const dirWidthLimit = 43
const totalWidthLimit = 80

// # Arguments
// - ancestorLines: Indicates whether a line should be drawn for the ancestor at each position.
// - start: Start position of the directory path.
func printDir(d dir, isLastChild bool, start int, ancestorLines []bool) (outputBytes []byte) {
	if !isLastChild {
		ancestorLines[start] = true
	}
	defer func() {
		ancestorLines[start] = false
	}()

	childStart := getChildStart(d.Path, start)

	// write parent prefix
	parentPrefix := getParentPrefix(ancestorLines)
	outputBytes = append(outputBytes, substr(string(parentPrefix), 0, start)...)

	// write self prefix
	selfPrefix := getSelfPrefix(isLastChild, start)
	outputBytes = append(outputBytes, selfPrefix...)

	// write path
	outputBytes = append(outputBytes, []byte(d.Path)...)

	// write space
	spaceWidth := dirWidthLimit - start - countUnicodeLength(string(selfPrefix)) - len(d.Path)
	outputBytes = append(outputBytes, bytes.Repeat([]byte(" "), spaceWidth)...)

	// write desc
	outputBytes = append(outputBytes, []byte("<= ")...)
	brokenDesc := make([]string, 0)
	for _, desc := range d.Desc {
		brokenDesc = append(brokenDesc, broke(desc, totalWidthLimit-dirWidthLimit)...)
	}

	for i, desc := range brokenDesc {
		if i == 0 {
			outputBytes = append(outputBytes, []byte(desc)...)
		} else {
			// write parent prefix
			parentPrefix := getParentPrefix(ancestorLines)
			outputBytes = append(outputBytes, substr(string(parentPrefix), 0, childStart)...)

			if len(d.Children) == 0 {
				outputBytes = append(outputBytes, []byte(" ")...)
			} else {
				outputBytes = append(outputBytes, []byte("│")...)
			}

			// write spaces
			spaceWidth := dirWidthLimit - childStart + 3 - 1 // +3 for "<= ", -1 for "│"
			outputBytes = append(outputBytes, bytes.Repeat([]byte(" "), spaceWidth)...)

			// write desc
			outputBytes = append(outputBytes, []byte(desc)...)
		}
		outputBytes = append(outputBytes, []byte("\n")...)
	}

	if len(brokenDesc) == 0 {
		outputBytes = append(outputBytes, []byte("[no description]\n")...)
	}

	for i, c := range d.Children {
		cIsFinalChild := i == len(d.Children)-1
		outputBytes = append(outputBytes, printDir(c, cIsFinalChild, childStart, ancestorLines)...)
	}

	return outputBytes
}

func broke(desc string, limit int) (broken []string) {
	if len(desc) <= limit {
		return []string{desc}
	}
	broken = make([]string, 0)

	lastSpaceIdx := 0
	start := 0
	for i := 0; i < len(desc); i++ {
		if desc[i] == ' ' {
			if i-start > limit {
				broken = append(broken, desc[start:lastSpaceIdx])
				start = lastSpaceIdx + 1
			}
			lastSpaceIdx = i
		}
	}
	broken = append(broken, desc[start:])
	return broken
}

func getParentPrefix(ancestorLines []bool) []byte {
	var parentPrefix []byte
	for i := 0; i < len(ancestorLines); i++ {
		if ancestorLines[i] {
			parentPrefix = append(parentPrefix, []byte("│")...)
		} else {
			parentPrefix = append(parentPrefix, []byte(" ")...)
		}
	}
	return parentPrefix
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

func countUnicodeLength(s string) (count int) {
	return len([]rune(s))
}
