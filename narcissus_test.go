package narcissus

import "os"

// Used in other tests
var wd, _ = os.Getwd()
var fakeroot = wd + "/fakeroot"

// Test structure
type foo struct {
	augeasPath string
	A          string `path:"a"`
}

type bar struct{}

type simpleValues struct {
	augeasPath string
	Str        string `path:"str"`
	Int        int    `path:"int"`
	Bool       bool   `path:"bool"`
}

type sliceValues struct {
	augeasPath string
	SlStr      []string   `path:"slstr"`
	SlInt      []int      `path:"slint"`
	SlBool     []bool     `path:"slbool"`
	SlStrSeq   []string   `type:"seq"`
	SlStruct   []mapEntry `path:"mapentry"`
}

type mapValues struct {
	augeasPath string
	Entries    map[string]mapEntry `path:"mstruct"`
	MStr       map[string]string   `path:"sub/*" key:"label"`
}

type mapEntry struct {
	Str   string   `path:"str"`
	Int   int      `path:"int"`
	Bool  bool     `path:"bool"`
	SlStr []string `path:"slstr"`
}

type noCapital struct {
	augeasPath string
	a          string `path:"a"`
}
