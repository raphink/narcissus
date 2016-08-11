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
	n := New(&aug)

	// Test one fstab
	entry := &FstabEntry{}
	err = n.Parse(entry, "/files/etc/fstab/1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if entry.File != "/" {
		t.Errorf("Expected file to be /, got %s", entry.File)
	}

	if len(entry.Opt) != 1 {
		t.Errorf("Expected one option, got %v", len(entry.Opt))
	}

	if entry.Opt["errors"].Value != "remount-ro" {
		t.Errorf("Expected option errors to have value remount-ro, got %s", entry.Opt["errors"].Value)
	}

	if entry.Passno != 1 {
		t.Errorf("Expected passno to be 1, got %v", entry.Passno)
	}

	// Test the whole fstab
	fstab := &Fstab{}
	err = n.Parse(fstab, "/files/etc/fstab")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(fstab.Entries) != 3 {
		t.Errorf("Expected 3 entries, got %v", len(fstab.Entries))
	}

	if fstab.Entries[0].File != "/" {
		t.Errorf("Expected file to be /, got %s", fstab.Entries[0].File)
	}

	if fstab.Entries[0].Opt["errors"].Value != "remount-ro" {
		t.Errorf("Expected option value to be remount-ro got %s", fstab.Entries[0].Opt["errors"].Value)
	}

	if len(fstab.Comments) != 6 {
		t.Errorf("Expected 6 comments, got %v", len(fstab.Comments))
	}

	if fstab.Comments[5].Comment != "/boot was on /dev/sda1 during installation" {
		t.Errorf("Expected comment text, got %s", fstab.Comments[5].Comment)
	}
}
