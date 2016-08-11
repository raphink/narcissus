package narcissus

import (
	"testing"

	"honnef.co/go/augeas"
)

type Fstab struct {
	Entries []FstabEntry
}

type FstabEntry struct {
	Spec    string `path:"spec"`
	File    string `path:"file"`
	Vfstype string `path:"vfstype"`
	Opt     []struct {
		Key   string `path:"."`
		Value string `path:"value"`
	} `path:"opt"`
	Dump   int `path:"dump"`
	Passno int `path:"passno"`
}

func TestFstab(t *testing.T) {

	aug, err := augeas.New("/home/raphink/go/src/github.com/raphink/narcissus/fakeroot", "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}

	// Test one fstab
	entry := &FstabEntry{}
	err = Parse(aug, entry, "/files/etc/fstab/1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if entry.File != "/" {
		t.Fatalf("Expected file to be /, got %s", entry.File)
	}

	/*
		if len(entry.Opt) != 1 {
			t.Fatalf("Expected one option, got %v", len(entry.Opt))
		}

		if entry.Opt[0].Key != "errors" {
			t.Fatalf("Expected option errors, got %s", entry.Opt[0].Key)
		}
		if entry.Opt[0].Value != "no-remounts" {
			t.Fatalf("Expected option value no-remounts, got %s", entry.Opt[0].Value)
		}
	*/

	if entry.Passno != 1 {
		t.Fatalf("Expected passno to be 1, got %v", entry.Passno)
	}

	// Test the whole fstab
	/*
		fstab := &Fstab{}
		err = Parse(aug, fstab, "/files/etc/fstab")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(fstab.Entries) != 3 {
			t.Fatalf("Expected 3 entries, got %v", len(fstab.Entries))
		}

		if fstab.Entries[0].File != "/" {
			t.Fatalf("Expected file to be /, got %s", fstab.Entries[0].File)
		}
	*/
}
