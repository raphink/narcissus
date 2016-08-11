package narcissus

import (
	"testing"

	"honnef.co/go/augeas"
)

func TestHosts(t *testing.T) {
	aug, err := augeas.New("/home/raphink/go/src/github.com/raphink/narcissus/fakeroot", "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Test one fstab
	host := &Host{}
	err = n.Parse(host, "/files/etc/hosts/1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if host.IPAddress != "127.0.0.1" {
		t.Errorf("Expected IP to be 127.0.0.1, got %s", host.IPAddress)
	}
}

func TestHost(t *testing.T) {
	aug, err := augeas.New("/home/raphink/go/src/github.com/raphink/narcissus/fakeroot", "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Test one fstab
	hosts, err := n.NewHosts()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(hosts.Hosts) != 5 {
		t.Errorf("Expected 3 hosts, got %v", len(hosts.Hosts))
	}

	if hosts.Hosts[2].IPAddress != "::1" {
		t.Errorf("Expected IP to be ::1, got %s", hosts.Hosts[2].IPAddress)
	}

	if len(hosts.Hosts[2].Aliases) != 2 {
		t.Errorf("Expected 2 aliases, got %v", len(hosts.Hosts[2].Aliases))
	}
}
