package narcissus

import (
	"fmt"
	"log"
	"os"
	"testing"

	"honnef.co/go/augeas"
)

func TestParseService(t *testing.T) {
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

func TestParseServices(t *testing.T) {
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

func TestWriteServices(t *testing.T) {
	// Use augeas.SaveNewFile once https://github.com/dominikh/go-augeas/issues/6 is fixed
	aug, err := augeas.New(fakeroot, "", 2)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Cleanup
	os.Remove(fakeroot + "/etc/services.augnew")

	services, err := n.NewServices()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	errStr, _ := aug.Get("/augeas//error/message")
	if errStr != "" {
		t.Errorf("Failed with %s", errStr)
	}

	err = wrapWrite(n, services, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	services.Services[100].Port = 42

	err = wrapWrite(n, services, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// check that file is changed
	expectedDiff := `--- orig
+++ new
@@ -113 +113 @@
-irc		194/tcp				# Internet Relay Chat
+irc		42/tcp				# Internet Relay Chat
`
	diff, err := diffNewContent("/etc/services")
	if err != nil {
		t.Errorf("Failed to compute diff: %v", err)
	} else if diff != expectedDiff {
		t.Errorf("Expected diff %s, got %s", expectedDiff, diff)
	}
}

func TestWriteServicesNewService(t *testing.T) {
	// Use augeas.SaveNewFile once https://github.com/dominikh/go-augeas/issues/6 is fixed
	aug, err := augeas.New(fakeroot, "", 2)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Cleanup
	os.Remove(fakeroot + "/etc/services.augnew")

	services, err := n.NewServices()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = wrapWrite(n, services, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	service := Service{
		Name:     "foo",
		Port:     12345,
		Protocol: "tcp",
		Comment:  "My foo service",
	}
	services.Services = append(services.Services, service)

	err = wrapWrite(n, services, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// check that file is changed
	expectedDiff := `--- orig
+++ new
@@ -612,0 +613 @@
+foo 12345/tcp # My foo service
`
	diff, err := diffNewContent("/etc/services")
	if err != nil {
		t.Errorf("Failed to compute diff: %v", err)
	} else if diff != expectedDiff {
		t.Errorf("Expected diff %s, got %s", expectedDiff, diff)
	}
}

func ExampleNarcissus_NewServices() {
	aug, err := augeas.New("/", "", augeas.None)
	if err != nil {
		log.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	services, err := n.NewServices()
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	fmt.Printf("Port=%v", services.Services[0].Port)
	// Output: Port=1
}

func ExampleNarcissus_NewService() {
	aug, err := augeas.New("/", "", augeas.None)
	if err != nil {
		log.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	service, err := n.NewService("ssh", "tcp")
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	fmt.Printf("Port=%v", service.Port)
	// Output: Port=22
}
