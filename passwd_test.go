package narcissus

import (
	"os"
	"testing"

	"honnef.co/go/augeas"
)

func TestParsePasswd(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	passwd, err := n.NewPasswd()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(passwd.Users) != 42 {
		t.Errorf("Expected 42 users, got %v", len(passwd.Users))
	}

	if passwd.Users["raphink"].UID != 1000 {
		t.Errorf("Expected user raphink to have uid 1000, got %v", passwd.Users["raphink"].UID)
	}
}

func TestParsePasswdUser(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Test one fstab
	user, err := n.NewPasswdUser("raphink")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.Account != "raphink" {
		t.Errorf("Expected account to be raphink, got %s", user.Account)
	}

	if user.UID != 1000 {
		t.Errorf("Expected uid to be 1000, got %v", user.UID)
	}

}

func TestWritePasswd(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.SaveNewFile)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Cleanup
	os.Remove(fakeroot + "/etc/passwd.augnew")

	passwd, err := n.NewPasswd()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = n.Write(passwd)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = aug.Save()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	errStr, _ := aug.Get("/augeas//error/message")
	if errStr != "" {
		t.Errorf("Failed with %s", errStr)
	}

	// check that file is unchanged
	if f, err := os.Stat(fakeroot + "/etc/passwd.augnew"); err == nil && f.Mode().IsRegular() {
		t.Errorf("Expected augnew file to be absent, was present")
	}

	var user = passwd.Users["raphink"]
	user.Shell = "/bin/sh"
	passwd.Users["raphink"] = user

	err = n.Write(passwd)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = aug.Save()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	errStr, _ = aug.Get("/augeas//error/message")
	if errStr != "" {
		t.Errorf("Failed with %s", errStr)
	}

	// check that file is changed
	expectedDiff := `--- orig
+++ new
@@ -41 +41 @@
-raphink:x:1000:1000:Raphaël Pinson,,,:/home/raphink:/bin/bash
+raphink:x:1000:1000:Raphaël Pinson,,,:/home/raphink:/bin/sh
`
	diff, err := diffNewContent("/etc/passwd")
	if err != nil {
		t.Errorf("Failed to compute diff: %v", err)
	} else if diff != expectedDiff {
		t.Errorf("Expected diff %s, got %s", expectedDiff, diff)
	}
}

func TestWritePasswdNewUser(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.SaveNewFile)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Cleanup
	os.Remove(fakeroot + "/etc/passwd.augnew")

	passwd, err := n.NewPasswd()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = n.Write(passwd)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = aug.Save()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	errStr, _ := aug.Get("/augeas//error/message")
	if errStr != "" {
		t.Errorf("Failed with %s", errStr)
	}

	// check that file is unchanged
	if f, err := os.Stat(fakeroot + "/etc/passwd.augnew"); err == nil && f.Mode().IsRegular() {
		t.Errorf("Expected augnew file to be absent, was present")
	}

	user, err := n.NewPasswdUser("foo")
	user.UID = 314
	user.GID = 314
	user.Name = "Foo Bar"
	user.Home = "/home/foo"
	user.Shell = "/bin/sh"
	user.Password = "XXX"
	passwd.Users["foo"] = *user

	err = n.Write(passwd)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = aug.Save()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	errStr, _ = aug.Get("/augeas//error/message")
	if errStr != "" {
		t.Errorf("Failed with %s", errStr)
	}

	// check that file is changed
	expectedDiff := `--- orig
+++ new
@@ -42,0 +43 @@
+foo:XXX:314:314:Foo Bar:/home/foo:/bin/sh
`
	diff, err := diffNewContent("/etc/passwd")
	if err != nil {
		t.Errorf("Failed to compute diff: %v", err)
	} else if diff != expectedDiff {
		t.Errorf("Expected diff %s, got %s", expectedDiff, diff)
	}
}
