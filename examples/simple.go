package main

import (
	"log"
	"strings"

	"github.com/raphink/narcissus"
	"honnef.co/go/augeas"
)

type group struct {
	augeasPath string
	Name       string   `path:"." value-from:"label"`
	Password   string   `path:"password"`
	GID        int      `path:"gid"`
	Users      []string `path:"user"`
}

func main() {
	aug, err := augeas.New("/", "", augeas.None)
	if err != nil {
		log.Fatal("Failed to create Augeas handler")
	}
	n := narcissus.New(&aug)

	user, err := n.NewPasswdUser("root")
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	log.Printf("UID=%v", user.UID)

	group := &group{
		augeasPath: "/files/etc/group/docker",
	}
	err = n.Parse(group)
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	log.Printf("GID=%v", group.GID)
	log.Printf("Users=%v", strings.Join(group.Users, ","))
}
