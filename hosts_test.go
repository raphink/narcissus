package narcissus

import (
	"fmt"
	"log"
	"testing"

	"honnef.co/go/augeas"
)

func TestHost(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	host := &Host{
		augeasPath: "/files/etc/hosts/1",
	}
	err = n.Parse(host)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if host.IPAddress != "127.0.0.1" {
		t.Errorf("Expected IP to be 127.0.0.1, got %s", host.IPAddress)
	}
}

func TestHosts(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	hosts, err := n.NewHosts()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(hosts.Hosts) != 5 {
		t.Errorf("Expected 5 hosts, got %v", len(hosts.Hosts))
	}

	if hosts.Hosts[2].IPAddress != "::1" {
		t.Errorf("Expected IP to be ::1, got %s", hosts.Hosts[2].IPAddress)
	}

	if len(hosts.Hosts[2].Aliases) != 2 {
		t.Errorf("Expected 2 aliases, got %v", len(hosts.Hosts[2].Aliases))
	}
}

func ExampleNarcissus_NewHosts() {
	aug, err := augeas.New(fakeroot, "", augeas.None)
	if err != nil {
		log.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	hosts, err := n.NewHosts()
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	fmt.Printf("IP=%v", hosts.Hosts[0].IPAddress)
	// Output: IP=127.0.0.1
}
