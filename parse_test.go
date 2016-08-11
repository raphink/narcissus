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
	Str        string `path:"str"`
	Int        int    `path:"int"`
	Bool       bool   `path:"bool"`
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
