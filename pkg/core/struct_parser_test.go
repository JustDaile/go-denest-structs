package core

import (
	"bytes"
	"strings"
	"testing"
)

func TestGetStructs(t *testing.T) {
	type testcase struct {
		def             []byte
		rootStructNames []string
	}
	tests := []testcase{
		{
			def: []byte(`
			type A struct {}
			type B struct {
				C struct {}
			}
			`),
			rootStructNames: []string{"A", "B"},
		},
		{
			def: []byte(`
			type A struct {}
			`),
			rootStructNames: []string{"A"},
		},
		{
			def: []byte(`
			type A struct {}
			type B struct {
				C struct {}
			}
			// conflicting name, parser doesn't care
			type B struct {
				C struct {}
			}
			`),
			rootStructNames: []string{"A", "B", "B"},
		},
	}
	for _, test := range tests {
		parser := StructParser{Buffer: bytes.NewBuffer(test.def)}

		_, defs := parser.GetStructs()
		if len(defs) != len(test.rootStructNames) {
			t.Errorf("expected definitions found %v to match total expected root structs name specified in testcase %v", len(defs), len(test.rootStructNames))
		}
		for i, def := range defs {
			if strings.Compare(def.GetName(), test.rootStructNames[i]) != 0 {
				t.Errorf("expected definition %v of %v to have name %s but was %s", i, len(tests), test.rootStructNames[i], def.GetName())
			}
		}
	}
}

func TestGetNamePosition(t *testing.T) {
	type testcase struct {
		def  []byte
		name string
	}
	tests := []testcase{
		{
			def:  []byte(`type A struct {}`),
			name: "A",
		},
		{
			def:  []byte(`A struct { B struct {} }`),
			name: "A",
		},
		{
			def:  []byte(`LongNameJustForTestingThatLongNameStructDoesntCauseFailureForUnknownReason struct { B struct {} }`),
			name: "LongNameJustForTestingThatLongNameStructDoesntCauseFailureForUnknownReason",
		},
		{
			def:  []byte(`Spaces      struct { B struct {} }`),
			name: "Spaces",
		},
		{
			def:  []byte(`type Spaces      struct { B struct {} }`),
			name: "Spaces",
		},
		{
			def:  []byte(`Tabbed	struct { B struct {} }`),
			name: "Tabbed",
		},
		{
			def:  []byte(`type Tabbed	struct { B struct {} }`),
			name: "Tabbed",
		},
		{
			def:  []byte(`type structInName	struct {}`),
			name: "structInName",
		},
		{
			def:  []byte(`type struct	struct {}`),
			name: "struct",
		},
		{
			def:  []byte(`struct_in_name	struct {}`),
			name: "struct_in_name",
		},
	}
	for _, test := range tests {
		def := StructDef(test.def)
		loc := def.getNamePosition()
		name := string(def[loc[0]:loc[1]])
		if strings.Compare(name, test.name) != 0 {
			t.Errorf("expected name position [%v:%v] to return name '%s' but was '%s'", loc[0], loc[1], test.name, name)
		}
	}
}

func TestRenameStruct(t *testing.T) {
	type testcase struct {
		def     []byte
		newName string
	}
	tests := []testcase{
		{
			def:     []byte(`type A struct {}`),
			newName: "B",
		},
		{
			def:     []byte(`A struct { B struct {} }`),
			newName: "C",
		},
	}
	for _, test := range tests {
		def := StructDef(test.def)
		def = def.RenameStruct(test.newName)
		if strings.Compare(def.GetName(), test.newName) != 0 {
			t.Errorf("expected struct with name %s to be renamed to %s", def.GetName(), test.newName)
		}
	}
}

func TestGetNestedStructLocation(t *testing.T) {
	type testcase struct {
		def             []byte
		matchStructName string
		offset          int
	}
	tests := []testcase{
		{
			def:             []byte(`A struct { B struct {} }`),
			matchStructName: "B",
		},
		{
			def:             []byte(`A struct { B struct {} C struct {} }`),
			offset:          1,
			matchStructName: "C",
		},
		{
			def:             []byte(`Tabbed	struct { B	struct {} C		struct {} }`),
			offset:          1,
			matchStructName: "C",
		},
	}
	for _, test := range tests {
		def := StructDef(test.def)
		locs := def.GetNestedStructLocations()
		if locs == nil {
			t.Fatalf("expected to find nested struct locations when using offset %v but got %v", test.offset, locs)
		}
		nested := StructDef(def[locs[test.offset][0]:locs[test.offset][1]])
		if strings.Compare(nested.GetName(), test.matchStructName) != 0 {
			t.Errorf("expected struct to get nested struct named %s when using offset %v but got this %s", test.matchStructName, test.offset, nested)
		}
	}
}
