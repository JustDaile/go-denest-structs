package core

import (
	"bytes"
	"strings"
	"testing"
)

func TestSeperateStruct(t *testing.T) {
	type testcase struct {
		bytes           []byte
		rootStructNames []string
	}
	testcases := []testcase{
		{
			// Should product structs [A, B] since B and C structs have identical struct bodies.
			// B should also be the first struct process, since processing read each struct as it is found.
			bytes:           []byte(`type A struct { B struct {} C struct {} }`),
			rootStructNames: []string{"A", "B"},
		},
		{
			// Should produce structs [A, B, D] since C and D structs have identical struct bodies.
			// D should also be the first struct processed, since processing reads each struct as it is found.
			bytes:           []byte(`type A struct { B struct { D struct {} } C struct {} }`),
			rootStructNames: []string{"A", "B", "D"},
		},
		{
			// Should produce structs [A, B, D, C, D_a46e2] since C and F structs have identical struct bodies.
			// D should also be the first struct processed, since processing reads each struct as it is found.
			// However since there is another struct with named D which does not match the existing definition for D,
			// The character '_' and the last 5 characters from the result of hashing the structs body is concatinated onto the structs name.
			// This is done wherever conflicts occur in the name of structs and ensures that no drastic name changes occur during processing.
			bytes:           []byte(`type A struct { B struct { D struct {} } C struct { D struct { F struct {} }} }`),
			rootStructNames: []string{"A", "B", "D", "C", "D_e1e02"},
		},
	}
	for _, test := range testcases {
		var buf *bytes.Buffer = bytes.NewBuffer([]byte{})
		if err := NewStructDenester(test.bytes).Process(buf); err != nil {
			t.Fatal(err)
		}
		parser := StructParser{Buffer: buf}
		_, structs := parser.GetStructs()
		if len(structs) != len(test.rootStructNames) {
			t.Errorf("expected structs returned %v to match the total struct names (%v) defined in testcase", len(structs), len(test.rootStructNames))
		}
		for _, strct := range structs {
			found := false
			for _, name := range test.rootStructNames {
				if strings.Compare(name, strct.GetName()) == 0 {
					found = true
				}
			}
			if !found {
				t.Errorf("unexpected struct called '%s' was found but was not expected within the testcase", strct.GetName())
			}
		}
	}
}
