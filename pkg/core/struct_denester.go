package core

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	utils "github.com/JustDaile/go-denest-structs/internal/pkg"
)

// StructDenester uses the StructParser to find structs defined within the given data and processes them into a new document
// while ensuring that no duplicate structs are added and that name conflicts are resolved.
// Conflicts within struct names are handled by hashing the body of a struct and appending the last 5 characters to then end of the original struct name.
type StructDenester struct {
	parser StructParser
	seen   map[string]string
}

// lookupHash private function is a somewhat rudimentary way of keeping track of which structs have already been processed and preventing duplication.
// The hash of the structs body, along with the name of the struct is stored so that conflicts can be found and fixed during processing.
func (denester StructDenester) lookupHash(hash string) *string {
	for h, n := range denester.seen {
		if strings.Compare(hash, h) == 0 {
			return &n
		}
	}
	return nil
}

// lookupName private function is a somewhat rudimentary way of keeping track of which structs have already been processed and preventing duplication.
// The name of the struct, along with the hash of the structs body is stored so that conflicts can be found and fixed during processing.
func (denester StructDenester) lookupName(name string) *string {
	for h, n := range denester.seen {
		if strings.Compare(name, n) == 0 {
			return &h
		}
	}
	return nil
}

// NewStructDenester factory to create a denester from an array of bytes.
func NewStructDenester(b []byte) *StructDenester {
	return &StructDenester{
		parser: StructParser{
			bytes.NewBuffer(b),
		},
	}
}

// Process seperates structs and writes out the resulting document without affecting any other content within within the original document and writes the result to the provided io.Writer.
func (denester StructDenester) Process(w io.Writer) error {
	denester.seen = map[string]string{}
	var bs StructDef = StructDef(denester.parser.Bytes()).Copy()

	locs, strcts := denester.parser.GetStructs()
	if len(locs) < 1 {
		return fmt.Errorf("unable to find any structs to process")
	}
	diff := 0
	for i, stct := range strcts {
		denestedStrcts := denester.seperateStructs(stct)
		var bout bytes.Buffer
		for _, denested := range denestedStrcts {
			bout.Write(denested)
			bout.Write([]byte("\n\n"))
		}
		bs = bs.overwrite(diff+locs[i][0], diff+locs[i][1], bout.Bytes())

		diff += bout.Len() - (locs[i][1] - locs[i][0])
	}
	w.Write(bs)
	return nil
}

// seperateStructs internal function is given a single struct definition and finds all nested/embedded structs and returns the results as defs.
func (denester StructDenester) seperateStructs(def StructDef) (defs []StructDef) {
	for {
		loc := def.GetNestedStructLocation()
		if loc == nil {
			break
		}
		nested := make(StructDef, (loc[1]-loc[0])+5)
		copy(nested, append([]byte("type "), def[loc[0]:loc[1]]...))
		arrLoc := utils.FindOpenAndCloseLocations(nested, '[', ']', 0, 0)
		if arrLoc != nil {
			nested = nested.overwrite(arrLoc[0], arrLoc[1], []byte{})
		}
		originalName := nested.GetName()
		hashedName := denester.lookupHash(nested.GetHash())
		namedHash := denester.lookupName(originalName)

		if namedHash == nil && hashedName == nil {
			hashedName = utils.Ptr(nested.GetName())
			namedHash = utils.Ptr(nested.GetHash())
			denester.seen[*namedHash] = *hashedName
			defs = append(defs, denester.seperateStructs(nested)...)
		} else if namedHash != nil && hashedName == nil {
			hashedName = utils.Ptr(fmt.Sprintf("%s_%s", nested.GetName(), nested.GetHash()[:5]))
			nested = nested.RenameStruct(*hashedName)
			denester.seen[nested.GetHash()] = *hashedName
			defs = append(defs, denester.seperateStructs(nested)...)
		}

		if arrLoc != nil {
			def = def.overwrite(loc[0], loc[1], []byte(fmt.Sprintf("%s []%s", originalName, *hashedName)))
		} else {
			def = def.overwrite(loc[0], loc[1], []byte(fmt.Sprintf("%s *%s", originalName, *hashedName)))
		}
	}
	defs = append([]StructDef{def}, defs...)
	return
}
