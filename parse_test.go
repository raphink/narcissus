package narcissus

import (
	"testing"

	"honnef.co/go/augeas"
)

type foo struct {
	augeasPath string
	A          string `path:"a"`
}

func TestParseNotAPtr(t *testing.T) {
	n := New(&augeas.Augeas{})
	err := n.Parse(foo{
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
	n := New(&augeas.Augeas{})
	f := "foo"
	err := n.Parse(&f)

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "not a struct" {
		t.Errorf("Expected error not a struct, got %s", err.Error())
	}
}

func TestParseFieldNotFound(t *testing.T) {
	n := New(&augeas.Augeas{})
	t.Skip("This causes a segfault with the Augeas lib. Open a bug!")
	err := n.Parse(&foo{
		augeasPath: "/files/some/path",
	})

	if err == nil {
		t.Error("Expected an error, got nothing")
	}
}
