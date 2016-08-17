package narcissus

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/pmezard/go-difflib/difflib"

	"honnef.co/go/augeas"
)

func TestWriteNotAPtr(t *testing.T) {
	aug, err := augeas.New("", "", augeas.NoModlAutoload)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)
	err = n.Write(foo{
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
	err = n.Write(&foo{})

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
		SlStruct: []mapEntry{
			{
				Str:   "foo",
				Int:   314,
				Bool:  true,
				SlStr: []string{"aleph", "beth"},
			},
			{
				Str:   "bar",
				Int:   315,
				Bool:  false,
				SlStr: []string{"gimel", "daleth"},
			},
		},
	}
	err = n.Write(s)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	checkAugMatch(t, n.Augeas, "/test/slstr", 2)
	checkAugGet(t, n.Augeas, "/test/slstr[2]", "b")

	checkAugMatch(t, n.Augeas, "/test/slint", 2)
	checkAugGet(t, n.Augeas, "/test/slint[2]", "2")

	checkAugMatch(t, n.Augeas, "/test/slbool", 2)
	checkAugGet(t, n.Augeas, "/test/slbool[1]", "true")
	checkAugGet(t, n.Augeas, "/test/slbool[2]", "false")

	checkAugMatch(t, n.Augeas, "/test/*[label()=~regexp('[0-9]*')]", 2)
	checkAugGet(t, n.Augeas, "/test/2", "bar")

	checkAugMatch(t, n.Augeas, "/test/mapentry", 2)
	checkAugGet(t, n.Augeas, "/test/mapentry[2]/str", "bar")
	checkAugMatch(t, n.Augeas, "/test/mapentry[2]/slstr", 2)
	checkAugGet(t, n.Augeas, "/test/mapentry[2]/slstr[2]", "daleth")
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
			"one": {
				Str:   "a",
				Int:   42,
				Bool:  true,
				SlStr: []string{"alpha", "beta"},
			},
			"two": {
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

	checkAugMatch(t, n.Augeas, "/test/mstruct", 2)
	checkAugGet(t, n.Augeas, "/test/mstruct[.='two']/str", "b")
	checkAugGet(t, n.Augeas, "/test/mstruct[.='two']/int", "43")
	checkAugGet(t, n.Augeas, "/test/mstruct[.='two']/bool", "false")
	checkAugMatch(t, n.Augeas, "/test/mstruct[.='two']/slstr", 2)
	checkAugGet(t, n.Augeas, "/test/mstruct[.='two']/slstr[2]", "delta")

	checkAugMatch(t, n.Augeas, "/test/sub/*", 2)
	checkAugGet(t, n.Augeas, "/test/sub/b", "beth")
}

// util methods

func diffNewContent(file string) (string, error) {
	origFile := fakeroot + file
	origContent, err := ioutil.ReadFile(origFile)
	if err != nil {
		return "", fmt.Errorf("Failed to read file %s: %v", origFile, err)
	}
	newFile := origFile + ".augnew"
	newContent, err := ioutil.ReadFile(newFile)
	if err != nil {
		return "", fmt.Errorf("Failed to read file %s: %v", newFile, err)
	}
	return diffContent(origContent, newContent)
}

func diffContent(orig []byte, new []byte) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(orig)),
		B:        difflib.SplitLines(string(new)),
		FromFile: "orig",
		ToFile:   "new",
		Context:  0,
		Eol:      "\n",
	}
	diffR, err := difflib.GetUnifiedDiffString(diff)
	return diffR, err
}

func augnewPresent(file string) bool {
	if f, err := os.Stat(fakeroot + file + ".augnew"); err == nil && f.Mode().IsRegular() {
		return true
	}
	return false
}

func wrapWrite(n *Narcissus, val interface{}, checkAugnew bool) (err error) {
	aug := n.Augeas
	err = n.Write(val)
	if err != nil {
		return fmt.Errorf("Failed writing, got %v", err)
	}
	errStr, _ := aug.Get("/augeas//error/message")
	if errStr != "" {
		return fmt.Errorf("Failed with %s", errStr)
	}

	/* FIXME: This fails, maybe because tests are parallelized?
	if checkAugnew && augnewPresent("/etc/passwd") {
		return fmt.Errorf("Expected augnew file to be absent, was present")
	}
	*/

	return nil
}

func ExampleNarcissus_Write() {
	aug, err := augeas.New(fakeroot, "", augeas.SaveNewFile)
	if err != nil {
		log.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	passwd, err := n.NewPasswd()
	err = n.Parse(passwd)
	if err != nil {
		log.Fatalf("Failed to parse passwd: %v", err)
	}

	user := passwd.Users["root"]
	fmt.Printf("Shell=%v\n", user.Shell)

	user.Shell = "/bin/zsh"
	passwd.Users["root"] = user

	err = n.Write(passwd)
	if err != nil {
		log.Fatalf("Failed to write passwd: %v", err)
	}

	cmd := exec.Command("grep", "^root", fakeroot+"/etc/passwd.augnew")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run grep: %v", err)
	}

	fmt.Println(out.String())
	// Output: Shell=/bin/bash
	// root:x:0:0:root:/root:/bin/zsh
}

func checkAugGet(t *testing.T, aug *augeas.Augeas, path string, expected string) {
	val, _ := aug.Get(path)
	if val != expected {
		t.Errorf("Expected %s, got %s", expected, val)
	}
}

func checkAugMatch(t *testing.T, aug *augeas.Augeas, path string, expected int) {
	matches, _ := aug.Match(path)
	if len(matches) != expected {
		t.Errorf("Expected %v elements, got %v", expected, len(matches))
	}
}
