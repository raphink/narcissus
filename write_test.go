package narcissus

import (
	"testing"

	"honnef.co/go/augeas"
)

type fooWrite struct {
	augeasPath string
	A          string `path:"a"`
}

func TestWriteNotAPtr(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	err = n.Write(fooWrite{
		augeasPath: "/files/some/path",
	})

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "invalid interface: not a ptr" {
		t.Errorf("Expected error not a ptr, got %s", err.Error())
	}
}

func TestWriteNotAStruct(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	f := "foo"
	err = n.Write(&f)

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "invalid interface: not a struct" {
		t.Errorf("Expected error not a struct, got %s", err.Error())
	}
}

func TestWriteNoAugeasPathValue(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	err = n.Write(&fooWrite{})

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "undefined path: no augeasPath value and no default" {
		t.Errorf("Expected error no augeasPath value and no default, got %s", err.Error())
	}
}

func TestWriteNoAugeasPathField(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	err = n.Write(&bar{})

	if err == nil {
		t.Error("Expected an error, got nothing")
	}

	if err.Error() != "undefined path: no augeasPath field" {
		t.Errorf("Expected error no augeasPath field, got %s", err.Error())
	}
}

func TestWriteSimpleField(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	s := &simpleValues{
		augeasPath: "/test",
		Str:        "foo",
		Int:        42,
		Bool:       true,
	}
	err = n.Write(s)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	got, _ := n.Augeas.Get("/test/str")
	if got != "foo" {
		t.Errorf("Expected foo, got %s", got)
	}

	got, _ = n.Augeas.Get("/test/int")
	if got != "42" {
		t.Errorf("Expected 42, got %v", got)
	}

	got, _ = n.Augeas.Get("/test/bool")
	if got != "true" {
		t.Errorf("Expected true, got %v", got)
	}
}
