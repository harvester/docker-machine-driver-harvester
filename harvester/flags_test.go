package harvester

import (
	"errors"
	"fmt"
	"testing"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
	digits   = "1234567890"
)

func TestParseLabels(t *testing.T) {
	type result struct {
		labels map[string]string
		errors error
	}

	tests := []struct {
		name        string
		input       string
		expectation result
	}{
		{
			name:        "empty labels string",
			input:       "",
			expectation: result{},
		},
		{
			name:  "just one label",
			input: "foobar=barfoo",
			expectation: result{
				labels: map[string]string{
					"foobar": "barfoo",
				},
			},
		},
		{
			name:  "multiple labels",
			input: "foo=bar,baz=bla,xxx=yyy",
			expectation: result{
				labels: map[string]string{
					"foo": "bar",
					"baz": "bla",
					"xxx": "yyy",
				},
			},
		},
		{
			name:  "multiple labels different order",
			input: "foo=bar,baz=bla,xxx=yyy",
			expectation: result{
				labels: map[string]string{
					"baz": "bla",
					"xxx": "yyy",
					"foo": "bar",
				},
			},
		},
		{
			name:  "long label value",
			input: fmt.Sprintf("foo=bar,baz=%s%s%s%s,xxx=yyy", alphabet, digits, alphabet, digits),
			expectation: result{
				labels: map[string]string{
					"foo": "bar",
					"baz": "hash_09h0hQ_z",
					"xxx": "yyy",
				},
			},
		},
		{
			name:  "error when key is missing value",
			input: "foo=bar,baz,xxx=yyy",
			expectation: result{
				labels: map[string]string{},
				errors: ParseLabelsSyntaxErr,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outLabels, outErrors := parseLabels(test.input)

			if test.expectation.errors == nil && outErrors != nil {
				t.Errorf("unexpected error: \"%s\"", outErrors)
			}

			if test.expectation.errors != nil && !errors.Is(outErrors, test.expectation.errors) {
				t.Errorf("error \"%s\" does not match expectation \"%s\"", outErrors, test.expectation.errors)
			}

			for olk, olv := range outLabels {
				elv, ok := test.expectation.labels[olk]

				if !ok {
					t.Errorf("unexpected label key: \"%s\"", olk)
				}

				if olv != elv {
					t.Errorf("unexpected label value \"%s\" for key \"%s\", expected \"%s\"", olv, olk, elv)
				}
			}

			for elk, elv := range test.expectation.labels {
				olv, ok := outLabels[elk]

				if !ok {
					t.Errorf("missing label key \"%s\"", elk)
				}

				if olv != elv {
					t.Errorf("label value \"%s\" for key \"%s\" does not match expected \"%s\"", olv, elk, elv)
				}
			}
		})
	}
}
