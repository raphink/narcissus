package narcissus

import "os"

// Used in other tests
var wd, _ = os.Getwd()
var fakeroot = wd + "/fakeroot"

// Test structure
type foo struct {
	augeasPath string
	A          string `narcissus:"a"`
}

type bar struct{}

type simpleValues struct {
	augeasPath string
	Str        string `narcissus:"str"`
	Int        int    `narcissus:"int"`
	Bool       bool   `narcissus:"bool"`
}

type sliceValues struct {
	augeasPath string
	SlStr      []string   `narcissus:"slstr"`
	SlInt      []int      `narcissus:"slint"`
	SlBool     []bool     `narcissus:"slbool"`
	SlStrSeq   []string   `narcissus:"seq"`
	SlStruct   []mapEntry `narcissus:"mapentry"`
}

type mapValues struct {
	augeasPath string
	Entries    map[string]mapEntry `narcissus:"mstruct"`
	MStr       map[string]string   `narcissus:"sub/*,key-from-label"`
}

type mapEntry struct {
	Str   string   `narcissus:"str"`
	Int   int      `narcissus:"int"`
	Bool  bool     `narcissus:"bool"`
	SlStr []string `narcissus:"slstr"`
}

type noCapital struct {
	augeasPath string
	a          string `narcissus:"a"`
}
