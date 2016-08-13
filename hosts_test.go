package narcissus

import (
	"fmt"
	"log"
	"os"
	"testing"

	"honnef.co/go/augeas"
)

func TestParseHost(t *testing.T) {
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

func TestParseHosts(t *testing.T) {
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

func TestWriteHosts(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.SaveNewFile)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Cleanup
	os.Remove(fakeroot + "/etc/hosts.augnew")

	hosts, err := n.NewHosts()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	errStr, _ := aug.Get("/augeas//error/message")
	if errStr != "" {
		t.Errorf("Failed with %s", errStr)
	}

	err = wrapWrite(n, hosts, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	hosts.Hosts[0].Canonical = "foo"

	err = wrapWrite(n, hosts, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// check that file is changed
	expectedDiff := `--- orig
+++ new
@@ -1 +1 @@
-127.0.0.1	localhost
+127.0.0.1	foo
`
	diff, err := diffNewContent("/etc/hosts")
	if err != nil {
		t.Errorf("Failed to compute diff: %v", err)
	} else if diff != expectedDiff {
		t.Errorf("Expected diff %s, got %s", expectedDiff, diff)
	}
}

func TestWriteHostsNewHost(t *testing.T) {
	aug, err := augeas.New(fakeroot, "", augeas.SaveNewFile)
	if err != nil {
		t.Fatal("Failed to create Augeas handler")
	}
	n := New(&aug)

	// Cleanup
	os.Remove(fakeroot + "/etc/hosts.augnew")

	hosts, err := n.NewHosts()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = wrapWrite(n, hosts, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	host := Host{
		IPAddress: "192.168.0.1",
		Canonical: "foo.example.com",
		Aliases:   []string{"foo", "bar"},
		Comment:   "Foo host",
	}
	hosts.Hosts = append(hosts.Hosts, host)

	err = wrapWrite(n, hosts, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// check that file is changed
	expectedDiff := `--- orig
+++ new
@@ -7,0 +8 @@
+192.168.0.1	foo.example.com foo bar # Foo host
`
	diff, err := diffNewContent("/etc/hosts")
	if err != nil {
		t.Errorf("Failed to compute diff: %v", err)
	} else if diff != expectedDiff {
		t.Errorf("Expected diff %s, got %s", expectedDiff, diff)
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
