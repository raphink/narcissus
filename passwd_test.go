package narcissus

import (
	"testing"

	"honnef.co/go/augeas"
)

func TestPasswd(t *testing.T) {
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

func TestPasswdUser(t *testing.T) {
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
