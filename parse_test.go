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
	aug, err := augeas.New("", "", augeas.None)
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

	if err.Error() != "not a ptr" {
		t.Errorf("Expected error not a ptr, got %s", err.Error())
	}
}

func TestParseNotAStruct(t *testing.T) {
	aug, err := augeas.New("", "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	f := "foo"
	err = n.Parse(&f)

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "not a struct" {
		t.Errorf("Expected error not a struct, got %s", err.Error())
	}
}

func TestParseFieldNotFound(t *testing.T) {
	aug, err := augeas.New("", "", augeas.None)
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
	aug, err := augeas.New("", "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	err = n.Parse(&foo{})

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "no augeasPath value and no default" {
		t.Errorf("Expected error no augeasPath value and no default, got %s", err.Error())
	}
}

func TestNoAugeasPathField(t *testing.T) {
	aug, err := augeas.New("", "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	err = n.Parse(&bar{})

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "no augeasPath field" {
		t.Errorf("Expected error no augeasPath field, got %s", err.Error())
	}
}

type simpleValues struct {
	augeasPath string
	Str        string   `path:"str"`
	Int        int      `path:"int"`
	Bool       bool     `path:"bool"`
	SlStr      []string `path:"slstr"`
	SlInt      []int    `path:"slint"`
}

func TestGetStringField(t *testing.T) {
	aug, err := augeas.New("", "", augeas.None)
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
}
