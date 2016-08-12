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

func TestWriteSliceField(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	s := &sliceValues{
		augeasPath: "/test",
		SlStr:      []string{"a", "b"},
		SlInt:      []int{1, 2},
		SlBool:     []bool{true, false},
		SlStrSeq:   []string{"foo", "bar"},
	}
	err = n.Write(s)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	matches, _ := n.Augeas.Match("/test/slstr")
	if len(matches) != 2 {
		t.Errorf("Expected 2 elements, got %v", len(matches))
	}
	got, _ := n.Augeas.Get("/test/slstr[2]")
	if got != "b" {
		t.Errorf("Expected element to be b, got %s", got)
	}

	matches, _ = n.Augeas.Match("/test/slint")
	if len(matches) != 2 {
		t.Errorf("Expected 2 elements, got %v", len(matches))
	}
	got, _ = n.Augeas.Get("/test/slint[2]")
	if got != "2" {
		t.Errorf("Expected element to be 2, got %v", got)
	}

	matches, _ = n.Augeas.Match("/test/slbool")
	if len(matches) != 2 {
		t.Errorf("Expected 2 elements, got %v", len(matches))
	}
	got, _ = n.Augeas.Get("/test/slbool[1]")
	if got != "true" {
		t.Errorf("Expected element to be true, got %v", got)
	}
	got, _ = n.Augeas.Get("/test/slbool[2]")
	if got != "false" {
		t.Errorf("Expected element to be false, got %v", got)
	}

	matches, _ = n.Augeas.Match("/test/*[label()=~regexp('[0-9]*')]")
	if len(matches) != 2 {
		t.Errorf("Expected 2 elements, got %v", len(matches))
	}
	got, _ = n.Augeas.Get("/test/2")
	if got != "bar" {
		t.Errorf("Expected element to be bar, got %s", got)
	}
}

func TestWriteMapField(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	m := &mapValues{
		augeasPath: "/test",
		Entries: map[string]mapEntry{
			"one": mapEntry{
				Str:   "a",
				Int:   42,
				Bool:  true,
				SlStr: []string{"alpha", "beta"},
			},
			"two": mapEntry{
				Str:   "b",
				Int:   43,
				Bool:  false,
				SlStr: []string{"gamma", "delta"},
			},
		},
		MStr: map[string]string{
			"a": "aleph",
			"b": "beth",
		},
	}
	err = n.Write(m)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	matches, _ := n.Augeas.Match("/test/mstruct")
	if len(matches) != 2 {
		t.Errorf("Expected 2 entries, got %v", len(matches))
	}
	got, _ := n.Augeas.Get("/test/mstruct[.='two']/str")
	if got != "b" {
		t.Errorf("Expected element to be b, got %s", got)
	}
	got, _ = n.Augeas.Get("/test/mstruct[.='two']/int")
	if got != "43" {
		t.Errorf("Expected element to be 43, got %s", got)
	}
	got, _ = n.Augeas.Get("/test/mstruct[.='two']/bool")
	if got != "false" {
		t.Errorf("Expected element to be false, got %s", got)
	}
	matches, _ = n.Augeas.Match("/test/mstruct[.='two']/slstr")
	if len(matches) != 2 {
		t.Errorf("Expected 2 entries, got %v", len(matches))
	}
	got, _ = n.Augeas.Get("/test/mstruct[.='two']/slstr[2]")
	if got != "delta" {
		t.Errorf("Expected element to be delta, got %v", got)
	}

	matches, _ = n.Augeas.Match("/test/sub/*")
	if len(matches) != 2 {
		t.Errorf("Expected 2 entries, got %v", len(matches))
	}
	got, _ = n.Augeas.Get("/test/sub/b")
	if got != "beth" {
		t.Errorf("Expected element to be beth, got %s", got)
	}
}
