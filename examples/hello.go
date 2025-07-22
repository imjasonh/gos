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
	color.Green("✓ Hello from gos example!")
	fmt.Printf("Arguments: %v\n", os.Args[1:])
}
