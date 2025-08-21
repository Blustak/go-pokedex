package main

import (
    "testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello  world   ",
			expected: []string{"hello", "world"},
		},
        {
            input: "",
            expected: []string{},
        },
        {
            input: "           ",
            expected: []string{},
        }, {
            input: `
            well, isn't this awkward.
            a multiline string.
            `,
            expected: []string{"well,", "isn't", "this", "awkward.","a","multiline","string."},
        },
	}
    for _, c := range cases {
        actual := cleanInput(c.input)
        if len(c.expected) != len(actual) {
            t.Errorf("error: function returned slice of length %d, expected %d", len(actual), len(c.expected))
        }

        for i := range actual {
            word := actual[i]
            expectedWord := c.expected[i]
            if word != expectedWord {
                t.Errorf("error: expected %s, got %s", expectedWord, word)
            }
        }
    }
}
