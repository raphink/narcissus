package narcissus

import (
	"testing"

	"honnef.co/go/augeas"
)

type foo struct {
	augeasPath string
	A          string `path:"a"`
}

type bar struct{}

func TestParseNotAPtr(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	err = n.Parse(foo{
		augeasPath: "/files/some/path",
	})

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "invalid interface: not a ptr" {
		t.Errorf("Expected error not a ptr, got %s", err.Error())
	}
}

func TestParseNotAStruct(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	f := "foo"
	err = n.Parse(&f)

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "invalid interface: not a struct" {
		t.Errorf("Expected error not a struct, got %s", err.Error())
	}
}

func TestParseFieldNotFound(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	err = n.Parse(&foo{
		augeasPath: "/files/some/path",
	})

	t.Skip("Fix this")

	if err == nil {
		t.Error("Expected an error, got nothing")
	}
}

func TestNoAugeasPathValue(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	err = n.Parse(&foo{})

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "undefined path: no augeasPath value and no default" {
		t.Errorf("Expected error no augeasPath value and no default, got %s", err.Error())
	}
}

func TestNoAugeasPathField(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	err = n.Parse(&bar{})

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "undefined path: no augeasPath field" {
		t.Errorf("Expected error no augeasPath field, got %s", err.Error())
	}
}

type simpleValues struct {
	augeasPath string
	Str        string `path:"str"`
	Int        int    `path:"int"`
	Bool       bool   `path:"bool"`
}

func TestGetSimpleField(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	n.Augeas.Set("/test/str", "foo")
	n.Augeas.Set("/test/int", "42")
	n.Augeas.Set("/test/bool", "true")
	n.Augeas.Set("/test/slstr[1]", "a")
	n.Augeas.Set("/test/slstr[2]", "b")
	n.Augeas.Set("/test/slint[1]", "1")
	n.Augeas.Set("/test/slint[2]", "2")
	n.Augeas.Set("/test/slbool[1]", "true")
	n.Augeas.Set("/test/slbool[2]", "false")
	s := &simpleValues{
		augeasPath: "/test",
	}
	err = n.Parse(s)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if s.Str != "foo" {
		t.Errorf("Expected foo, got %s", s.Str)
	}

	if s.Int != 42 {
		t.Errorf("Expected 42, got %v", s.Int)
	}

	if s.Bool != true {
		t.Errorf("Expected true, got %v", s.Bool)
	}
}

func TestGetSimpleFieldWrongTypes(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	n.Augeas.Set("/test/int", "a")
	s := &simpleValues{
		augeasPath: "/test",
	}
	err = n.Parse(s)

	if err == nil {
		t.Error("Expected an error, got nil")
	}

	if err.Error() != "failed to retrieve field Int: failed to convert a to int: strconv.ParseInt: parsing \"a\": invalid syntax" {
		t.Errorf("Expected int conversion error, got %v", err)
	}

	n.Augeas.Remove("/test/int")
	n.Augeas.Set("/test/bool", "a")
	err = n.Parse(s)

	if err == nil {
		t.Error("Expected an error, got nil")
	}

	if err.Error() != "failed to retrieve field Bool: failed to convert a to bool: strconv.ParseBool: parsing \"a\": invalid syntax" {
		t.Errorf("Expected bool conversion error, got %v", err)
	}
}

type sliceValues struct {
	augeasPath string
	SlStr      []string `path:"slstr"`
	SlInt      []int    `path:"slint"`
	SlBool     []bool   `path:"slbool"`
	SlStrSeq   []string `type:"seq"`
}

func TestGetSliceField(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	n.Augeas.Set("/test/slstr[1]", "a")
	n.Augeas.Set("/test/slstr[2]", "b")
	n.Augeas.Set("/test/slint[1]", "1")
	n.Augeas.Set("/test/slint[2]", "2")
	n.Augeas.Set("/test/slbool[1]", "true")
	n.Augeas.Set("/test/slbool[2]", "false")
	n.Augeas.Set("/test/1", "foo")
	n.Augeas.Set("/test/2", "bar")
	s := &sliceValues{
		augeasPath: "/test",
	}
	err = n.Parse(s)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(s.SlStr) != 2 {
		t.Errorf("Expected 2 elements, got %v", len(s.SlStr))
	}
	if s.SlStr[1] != "b" {
		t.Errorf("Expected element to be b, got %s", s.SlStr[1])
	}

	if len(s.SlInt) != 2 {
		t.Errorf("Expected 2 elements, got %v", len(s.SlInt))
	}
	if s.SlInt[1] != 2 {
		t.Errorf("Expected element to be 2, got %v", s.SlInt[1])
	}

	if len(s.SlBool) != 2 {
		t.Errorf("Expected 2 elements, got %v", len(s.SlBool))
	}
	if s.SlBool[0] != true {
		t.Errorf("Expected element to be true, got %v", s.SlBool[0])
	}
	if s.SlBool[1] != false {
		t.Errorf("Expected element to be false, got %v", s.SlBool[1])
	}

	if len(s.SlStrSeq) != 2 {
		t.Errorf("Expected 2 elements, got %v", len(s.SlStrSeq))
	}
	if s.SlStrSeq[1] != "bar" {
		t.Errorf("Expected element to be bar, got %s", s.SlStrSeq[1])
	}
}

type mapValues struct {
	augeasPath string
	Entries    map[string]struct {
		Str   string   `path:"str"`
		Int   int      `path:"int"`
		Bool  bool     `path:"bool"`
		SlStr []string `path:"slstr"`
	} `path:"mstruct"`
}

func TestGetMapField(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	n.Augeas.Set("/test/mstruct[1]", "one")
	n.Augeas.Set("/test/mstruct[1]/str", "a")
	n.Augeas.Set("/test/mstruct[1]/int", "42")
	n.Augeas.Set("/test/mstruct[1]/bool", "true")
	n.Augeas.Set("/test/mstruct[1]/slstr[1]", "alpha")
	n.Augeas.Set("/test/mstruct[1]/slstr[2]", "beta")
	n.Augeas.Set("/test/mstruct[2]", "two")
	n.Augeas.Set("/test/mstruct[2]/str", "b")
	n.Augeas.Set("/test/mstruct[2]/int", "43")
	n.Augeas.Set("/test/mstruct[2]/bool", "false")
	n.Augeas.Set("/test/mstruct[2]/slstr[1]", "gamma")
	n.Augeas.Set("/test/mstruct[2]/slstr[2]", "delta")
	m := &mapValues{
		augeasPath: "/test",
	}
	err = n.Parse(m)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(m.Entries) != 2 {
		t.Errorf("Expected 2 entries, got %v", len(m.Entries))
	}
	if m.Entries["two"].Str != "b" {
		t.Errorf("Expected element to be b, got %s", m.Entries["two"].Str)
	}
}

type noCapital struct {
	augeasPath string
	a          string `path:"a"`
}

func TestSetField(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	n.Augeas.Set("/test/a", "a")
	s := &noCapital{
		augeasPath: "/test",
	}
	err = n.Parse(s)

	if err == nil {
		t.Error("Expected an error, got nil")
	}

	if err.Error() != "cannot set field a" {
		t.Errorf("Expected setField error, got %v", err)
	}
}
