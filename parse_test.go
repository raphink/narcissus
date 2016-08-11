package narcissus

import (
	"testing"

	"honnef.co/go/augeas"
)

type foo struct {
}

func TestParseNotAPtr(t *testing.T) {
	n := New(&augeas.Augeas{})
	err := n.Parse(foo{}, "/files/some/path")

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "not a ptr" {
		t.Errorf("Expected error not a ptr, got %s", err.Error())
	}
}

func TestParseNotAStruct(t *testing.T) {
	n := New(&augeas.Augeas{})
	f := "foo"
	err := n.Parse(&f, "/files/some/path")

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "not a struct" {
		t.Errorf("Expected error not a struct, got %s", err.Error())
	}
}
