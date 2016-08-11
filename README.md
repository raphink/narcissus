Narcissus
=========

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/raphink/narcissus)
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
		log.Fatalf("Failed to retrieve user: %v" err)
	}

	log.Printf("UID=%v", user.UID)
}
```

## Available types

### Fstab

* [`Fstab`](https://godoc.org/github.com/raphink/narcissus#Fstab) maps a whole `/etc/fstab` file
* [`FstabEntry`](https://godoc.org/github.com/raphink/narcissus#FstabEntry) maps a single `/etc/fstab` entry

### Hosts

* [`Hosts`](https://godoc.org/github.com/raphink/narcissus#Hosts) maps a whole `/etc/hosts` file
* [`Host`](https://godoc.org/github.com/raphink/narcissus#Host) maps a single `/etc/hosts` entry

### Passwd

* [`Passwd`](https://godoc.org/github.com/raphink/narcissus#Passwd) maps a whole `/etc/passwd` file
* [`PasswdUser`](https://godoc.org/github.com/raphink/narcissus#PasswdUser) maps a single `/etc/passwd` entry

### Services

* [`Services`](https://godoc.org/github.com/raphink/narcissus#Services) maps a whole `/etc/services` file
* [`Service`](https://godoc.org/github.com/raphink/narcissus#Service) maps a single `/etc/services` entry


## Mapping your own structures


```go
import (
	"log"

	"honnef.co/go/augeas"
	"github.com/raphink/narcissus"
)

type group struct {
	augeasPath string
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

	group := &group{
		augeasPath: "/files/etc/group/docker",
	}
	err = n.Parse(group)
	if err != nil {
		log.Fatalf("Failed to retrieve group: %v", err)
	}

	log.Printf("GID=%v", group.GID)
}
```
