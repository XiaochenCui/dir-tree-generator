# dir-tree-generator

Generate directory tree with description from yaml file.

## Usage

### As a Binary

```bash
go install github.com/XiaochenCui/dir_tree_generator@latest
dir_tree_generator input.yaml
```

### As a Library

```go
package main

import (
    "os"
    "log"

    "github.com/XiaochenCui/dir_tree_generator"
)

func main() {
    filepath := "input.yaml"
    yamlBytes, err := os.ReadFile(filepath)
    if err != nil {
        log.Fatal(err)
        return
    }

    output, err := os.ReadFile(wantOutputPath)
    if err != nil {
        log.Fatal(err)
        return
    }

    log.Println(string(output))
}
```

## Example

input:

![input](/image/input-1.png)

output:

![output](/image/output-1.png)

- The child directories always start from the right position.
- The description is wrapped correctly.
