package main

import (
	"log"

	"github.com/raphink/narcissus"
	"honnef.co/go/augeas"
)

type group struct {
	Name     string `path:"." value-from:"label"`
	Password string `path:"password"`
	GID      int    `path:"gid"`
}

func main() {
	aug, err := augeas.New("/", "", augeas.None)
	if err != nil {
		log.Fatal("Failed to create Augeas handler")
	}
	n := narcissus.New(&aug)

	user := &narcissus.PasswdUser{}
	err = n.Parse(user, "/files/etc/passwd/root")
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	log.Printf("UID=%v", user.Shell)

	group := &group{}
	err = n.Parse(group, "/files/etc/group/docker")
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	log.Printf("GID=%v", group.GID)
}
