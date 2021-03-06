package narcissus

import (
	"fmt"
	"log"
	"os"
	"testing"

	"honnef.co/go/augeas"
)

func TestParseFstabEntry(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Test one fstab
	entry := &FstabEntry{
		augeasPath: "/files/etc/fstab/1",
	}
	err = n.Parse(entry)
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
}

func TestParseFstab(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	fstab, err := n.NewFstab()
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

func TestWriteFstab(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.SaveNewFile)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Cleanup
	os.Remove(fakeroot + "/etc/fstab.augnew")

	fstab, err := n.NewFstab()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	errStr, _ := aug.Get("/augeas//error/message")
	if errStr != "" {
		t.Errorf("Failed with %s", errStr)
	}

	err = wrapWrite(n, fstab, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	fstab.Entries[0].File = "/foo"

	err = wrapWrite(n, fstab, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// check that file is changed
	expectedDiff := `--- orig
+++ new
@@ -8 +8 @@
-/dev/mapper/wrk8--vg-root /               ext4    errors=remount-ro 0       1
+/dev/mapper/wrk8--vg-root /foo               ext4    errors=remount-ro 0       1
`
	diff, err := diffNewContent("/etc/fstab")
	if err != nil {
		t.Errorf("Failed to compute diff: %v", err)
	} else if diff != expectedDiff {
		t.Errorf("Expected diff %s, got %s", expectedDiff, diff)
	}
}

func TestWriteFstabNewEntry(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.SaveNewFile)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Cleanup
	os.Remove(fakeroot + "/etc/fstab.augnew")

	fstab, err := n.NewFstab()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = wrapWrite(n, fstab, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	tab := FstabEntry{
		Spec:    "/dev/foo",
		File:    "/foo",
		Vfstype: "ext4",
		Opt: map[string]FstabOpt{
			"defaults": {},
		},
		Dump:   0,
		Passno: 0,
	}
	fstab.Entries = append(fstab.Entries, tab)

	err = wrapWrite(n, fstab, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// check that file is changed
	expectedDiff := `--- orig
+++ new
@@ -11,0 +12 @@
+/dev/foo	/foo	ext4	defaults	0 0
`
	diff, err := diffNewContent("/etc/fstab")
	if err != nil {
		t.Errorf("Failed to compute diff: %v", err)
	} else if diff != expectedDiff {
		t.Errorf("Expected diff %s, got %s", expectedDiff, diff)
	}
}

func TestWriteFstabDeleteEntry(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.SaveNewFile)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Cleanup
	os.Remove(fakeroot + "/etc/fstab.augnew")

	fstab, err := n.NewFstab()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = wrapWrite(n, fstab, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Remove fstab.Entries[1]
	fstab.Entries = append(fstab.Entries[:1], fstab.Entries[2:]...)

	t.Skip("Slice purge not implemented")

	err = wrapWrite(n, fstab, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// check that file is changed
	expectedDiff := `--- orig
+++ new
@@ -11,0 +12 @@
+/dev/foo	/foo	ext4	defaults	0 0
`
	diff, err := diffNewContent("/etc/fstab")
	if err != nil {
		t.Errorf("Failed to compute diff: %v", err)
	} else if diff != expectedDiff {
		t.Errorf("Expected diff %s, got %s", expectedDiff, diff)
	}
}

func ExampleNarcissus_NewFstab() {
	aug, err := augeas.New(fakeroot, "", augeas.None)
	if err != nil {
		log.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	fstab, err := n.NewFstab()
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	fmt.Printf("File=%v", fstab.Entries[0].File)
	// Output: File=/
}
