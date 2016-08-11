package narcissus

import (
	"testing"

	"honnef.co/go/augeas"
)

func TestPasswd(t *testing.T) {

	aug, err := augeas.New("/home/raphink/go/src/github.com/raphink/narcissus/fakeroot", "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Test one fstab
	user := &PasswdUser{}
	err = n.Parse(user, "/files/etc/passwd/raphink")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.Account != "raphink" {
		t.Errorf("Expected account to be raphink, got %s", user.Account)
	}

	if user.Uid != 1000 {
		t.Errorf("Expected uid to be 1000, got %v", user.Uid)
	}

	// Test full file
	passwd := &Passwd{}
	err = n.Parse(passwd, "/files/etc/passwd")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(passwd.Users) != 42 {
		t.Errorf("Expected 42 users, got %v", len(passwd.Users))
	}

	if passwd.Users["raphink"].Uid != 1000 {
		t.Errorf("Expected user raphink to have uid 1000, got %v", passwd.Users["raphink"].Uid)
	}
}
