package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"regexp"

	utils "github.com/JustDaile/go-denest-structs/internal/pkg"
)

var (
	structSigRegxp      *regexp.Regexp
	anonyStructSigRegxp *regexp.Regexp
)

func init() {
	structSigRegxp, _ = regexp.Compile(`type\s+[A-Za-z_0-9]+\s+struct`)
	anonyStructSigRegxp, _ = regexp.Compile(`([a-zA-Z_][a-zA-Z0-9_]*)\s*(\[\])?\s*struct\s+{`)
}

// StructParser uses regex and loops to read through a byte.Buffer and parse out structs and their respective parts (name, body, embedded structs).
// Other support could be added such as parsing and adding tags.
type StructParser struct {
	*bytes.Buffer
}

type StructDef []byte

// getNamePosition gets the start and end of the first seen structs name as an array []int{start, end} or nil if no struct found.
func (s StructDef) getNamePosition() []int {
	ex := bytes.Index(s, []byte("struct {")) - 1
	if ex == -1 {
		return nil
	}
	var lastByte byte = s[ex]
	for (lastByte == ' ' || lastByte == '	' || lastByte == '[' || lastByte == ']') && ex > 0 {
		ex--
		lastByte = s[ex]
	}
	sx := ex - 1
	for lastByte != ' ' && sx > 0 {
		lastByte = s[sx]
		if lastByte != ' ' && lastByte != '	' {
			sx--
		}
	}
	if sx != 0 {
		sx += 1
	}
	return []int{sx, ex + 1}
}

// Copy returns a copy of the struct def as []byte
func (s StructDef) Copy() (b []byte) {
	b = make([]byte, len(s))
	copy(b, s)
	return
}

// GetName returns the name of the first found struct within the StructDef.
func (s StructDef) GetName() string {
	loc := s.getNamePosition()
	d := s[loc[0]:loc[1]]
	b := make([]byte, len(d))
	copy(b, d)
	return string(d)
}

// RenameStruct returns a Copy of the StructDef with the first found structs name replaced with the 'name' provided.
func (s StructDef) RenameStruct(name string) StructDef {
	loc := s.getNamePosition()
	return s.overwrite(loc[0], loc[1], []byte(name))
}

// overwrite package function returns the result of overwriting the bytes between StructDef[:sx] and StructDef[ex:] with the provide b ([]byte)
func (s StructDef) overwrite(sx, ex int, b []byte) (def StructDef) {
	def = make([]byte, len(s)+len(b))
	def = append(s[:sx], append(b, s[ex:]...)...)
	return
}

// GetHash returns a hash based on the bytes within the StructDef.
// Should be used to prevent duplicate structs.
func (s StructDef) GetHash() string {
	h := sha256.New()
	bs := utils.GetBytesBetween(s, '{', '}', 0, 0)
	c := []byte{}
	for _, b := range bs {
		if b != '\n' && b != '\t' && b != ' ' {
			c = append(c, b)
		}
	}
	h.Write(c)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetBody returns the first body found within the struct def (bytes between '{' and '}' inclusive).
func (s StructDef) GetBody() (b []byte) {
	bodyBytes := utils.GetBytesBetween(s, '{', '}', 0, 0)
	b = make([]byte, len(bodyBytes))
	copy(b, bodyBytes)
	return
}

// GetStructs - Find all none nested/embedded structs
func (parser StructParser) GetStructs() (locs [][]int, defs []StructDef) {
	bs := parser.Bytes()
	loc := structSigRegxp.FindAllIndex(bs, -1)
	if loc == nil {
		return
	}
	for _, se := range loc {
		var b []byte = make([]byte, (se[1]+1)-se[0])
		copy(b, bs[se[0]:se[1]+1])
		stct := append(b, utils.GetBytesBetween(bs, '{', '}', se[0], se[1])...)
		defs = append(defs, stct)
		locs = append(locs, []int{se[0], se[0] + len(stct)})
	}
	return
}

// GetNestedStructLocations returns the start and end position of each nested/embedded struct within the StructDef as an [][]int or array of int[2]{start, end}
func (s StructDef) GetNestedStructLocations() (locs [][]int) {
	bodyOffsets := utils.FindOpenAndCloseLocations(s, '{', '}', 0, 0)
	iOffset := bodyOffsets[0]
	for {
		if iOffset >= len(s) {
			return
		}
		loc := anonyStructSigRegxp.FindIndex(s[iOffset:])
		if loc == nil {
			return
		}
		bodyOffset := iOffset + loc[0]
		oc := utils.FindOpenAndCloseLocations(s, '{', '}', bodyOffset, bodyOffset)
		locs = append(locs, []int{bodyOffset, oc[1]})
		iOffset = oc[1]
	}
}

// GetNestedStructLocations returns the start and end position of the next nested/embedded struct within the StructDef as int[]{start, end}
func (s StructDef) GetNestedStructLocation() []int {
	bodyOffsets := utils.FindOpenAndCloseLocations(s, '{', '}', 0, 0)
	loc := anonyStructSigRegxp.FindIndex(s[bodyOffsets[0]:])
	if loc == nil {
		return nil
	}
	bodyOffset := bodyOffsets[0] + loc[0]
	oc := utils.FindOpenAndCloseLocations(s, '{', '}', bodyOffset, bodyOffset)
	return []int{bodyOffset, oc[1]}
}
