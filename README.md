Narcissus
=========

[![Documentation](https://img.shields.io/badge/godoc-narcissus-blue.svg)](https://godoc.org/github.com/raphink/narcissus)
[![Build Status](https://img.shields.io/travis/raphink/narcissus/master.svg)](https://travis-ci.org/raphink/narcissus)
[![Coverage Status](https://img.shields.io/coveralls/raphink/narcissus.svg)](https://coveralls.io/r/raphink/narcissus?branch=master)

This go package aims to provide reflection for the Augeas library.

## Example

```go
import (
	"log"

	"honnef.co/go/augeas"
	"github.com/raphink/narcissus"
)

func main() {
	aug, err := augeas.New("/", "", augeas.None)
	if err != nil {
		log.Fatal("Failed to create Augeas handler")
	}
  n := narcissus.New(&aug)

	user := n.NewPasswdUser("raphink")
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	log.Printf("UID=%v", user.UID)
}
```



## Mapping your own structures


```go
import (
	"log"

	"honnef.co/go/augeas"
	"github.com/raphink/narcissus"
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

	group := &group{}
	err = n.Parse(group, "/files/etc/group/docker")
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	log.Printf("GID=%v", group.GID)
}
```
