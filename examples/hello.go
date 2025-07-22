#!/usr/bin/env gos run
// /// script
// dependencies = [
//     "github.com/fatih/color@v1.18.0",
//	   "github.com/stretchr/testify@v1.10.0",
// ]
// ///

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func main() {
	color.Green("✓ Hello from gos example!")
	fmt.Printf("Arguments: %v\n", os.Args[1:])
}

func TestFoo(t *testing.T) {
	assert.Equal(t, 1, 1, "This should always pass")
}
