#!/usr/bin/env gos test
// /// script
// dependencies = [
//     "github.com/stretchr/testify@v1.10.0",
// ]
// ///

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Add(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}

func TestAdd(t *testing.T) {
	assert.Equal(t, 5, Add(2, 3), "2 + 3 should equal 5")
	assert.Equal(t, 0, Add(-1, 1), "-1 + 1 should equal 0")
	assert.Equal(t, -5, Add(-2, -3), "-2 + -3 should equal -5")
}

func TestMultiply(t *testing.T) {
	require.Equal(t, 6, Multiply(2, 3), "2 * 3 should equal 6")
	require.Equal(t, 0, Multiply(0, 100), "0 * 100 should equal 0")
	require.Equal(t, -6, Multiply(-2, 3), "-2 * 3 should equal -6")
}

func TestTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 5, 3, 8},
		{"negative numbers", -5, -3, -8},
		{"mixed signs", -5, 3, -2},
		{"with zero", 0, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}