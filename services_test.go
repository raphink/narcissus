package narcissus

import (
	"testing"

	"honnef.co/go/augeas"
)

func TestService(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Test one service
	service, err := n.NewService("ssh", "tcp")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if service.Port != 22 {
		t.Errorf("Expected port 22, got %v", service.Port)
	}

	if service.Comment != "SSH Remote Login Protocol" {
		t.Errorf("Expected SSH comment, got %s", service.Comment)
	}
}

func TestServices(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.None)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	services, err := n.NewServices()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(services.Services) != 557 {
		t.Errorf("Expected 557 services, got %v", len(services.Services))
	}
}
