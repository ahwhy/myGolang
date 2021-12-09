package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	// assert equality
	assert.Equal(t, Add(5), 7, "they should be equal")
}

func TestCal(t *testing.T) {
	ass := assert.New(t)
	var tests = []struct {
		input    int
		expected int
	}{
		{2, 4},
		{-1, 1},
		{0, 2},
		{-5, -3},
		{999999997, 999999999},
	}
	for _, test := range tests {
		ass.Equal(Add(test.input), test.expected)
	}
}
