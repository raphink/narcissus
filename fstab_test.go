package narcissus

import (
	"testing"

	"honnef.co/go/augeas"
)

func TestFstab(t *testing.T) {

	aug, err := augeas.New("/home/raphink/go/src/github.com/raphink/narcissus/fakeroot", "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}

	n := &Narcissus{
		Augeas: &aug,
	}

	// Test one fstab
	entry := &FstabEntry{}
	err = n.Parse(entry, "/files/etc/fstab/1")
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
		err = n.Parse(fstab, "/files/etc/fstab")
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
