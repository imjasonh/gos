# gos - Go Script Runner

A tool that enables running Go files as scripts with inline dependency declarations, inspired by Python's inline script dependencies ([PEP 723](https://peps.python.org/pep-0723/)), uv's [`run` command](https://docs.astral.sh/uv/guides/scripts/), and this [Hacker News discussion](https://news.ycombinator.com/item?id=44641746).

![Demo](demo.gif)

## Motivation

Go is a compiled language that typically requires a full project structure with `go.mod` files for dependency management. This tool allows you to write self-contained Go scripts with dependencies declared inline, making Go more suitable for quick scripts and automation tasks.

## Features

- Run Go files directly with inline dependencies
- No need for a separate `go.mod` file
- Automatic dependency resolution via `go mod tidy`
- Shebang support for executable scripts
- Temporary build environment (doesn't pollute your workspace)
- Support for both running scripts and testing them

## Installation

```bash
go install github.com/imjasonh/gos
```

## Usage

### Basic Usage

```bash
# Run a Go script
gos run script.go [args...]

# Run tests in a Go script
gos test script.go [args...]
```

### With Shebang

Add this to the top of your Go file:
```go
#!/usr/bin/env gos run
```

For test files:
```go
#!/usr/bin/env gos test
```

Then make it executable:
```bash
chmod +x script.go
./script.go [args...]
```

**Note:** Adding a shebang makes the file invalid Go syntax and may break some tools like `gofmt` and IDE formatting. The shebang is optional - you can always use `gos run script.go` instead.

## Dependency Format

Dependencies are declared in a special comment block at the top of your Go file:

```go
#!/usr/bin/env gos run
// /// script
// dependencies = [
//     "github.com/fatih/color@v1.18.0",
//     "github.com/spf13/cobra@v1.8.0",
// ]
// ///

package main

import (
    "github.com/fatih/color"
    // ... your imports
)

func main() {
    color.Green("Hello from gos!")
}
```

## Examples

### Simple Script with Colors

```go
#!/usr/bin/env gos run
// /// script
// dependencies = [
//     "github.com/fatih/color@v1.18.0",
// ]
// ///

package main

import (
    "fmt"
    "os"
    "github.com/fatih/color"
)

func main() {
    color.Green("✓ Hello from gos!")
    fmt.Printf("Arguments: %v\n", os.Args[1:])
}
```

### Combined Script and Tests

You can have both a runnable script and tests in the same file:

```go
#!/usr/bin/env gos run
// /// script
// dependencies = [
//     "github.com/stretchr/testify@v1.10.0",
// ]
// ///

package main

import (
    "fmt"
    "testing"
    "github.com/stretchr/testify/assert"
)

func main() {
    fmt.Println("Running main function")
}

func TestSomething(t *testing.T) {
    assert.Equal(t, 1, 1)
}
```

Then:
- `gos run script.go` or `./script.go` - runs the main function
- `gos test script.go` - runs the tests (automatically renames to `*_test.go`)

## How It Works

1. `gos` parses the special comment block to extract dependencies
2. Creates a temporary directory with a generated `go.mod`
3. Runs `go mod tidy` to fetch dependencies
4. For `run`: builds and executes your script
5. For `test`: renames to `*_test.go` if needed and runs `go test`
6. Cleans up the temporary directory

## Comparison to Other Tools

- **Python + uv**: Python's PEP 723 allows inline script dependencies, and uv's `run` command executes them
- **Deno**: Supports URL imports and can run TypeScript directly
- **Bun**: Can run TypeScript files with automatic dependency installation
- **gos**: Brings similar convenience to Go while maintaining Go's type safety and compilation

## Limitations

- Dependencies must be explicitly versioned or use `latest`
- No caching of built binaries (rebuilds each time)
- Requires `go` to be installed and available in PATH
- Test files must contain valid test functions for `gos test` to work

## Future Improvements

- Binary caching for faster repeated runs
- Support for `go.sum` verification
- Additional commands beyond `run`
- Dependency caching across scripts
- Support for replace directives
