Narcissus
=========

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/raphink/narcissus)
[![Build Status](https://img.shields.io/travis/raphink/narcissus/master.svg)](https://travis-ci.org/raphink/narcissus)
[![Coverage Status](https://img.shields.io/coveralls/raphink/narcissus.svg)](https://coveralls.io/r/raphink/narcissus?branch=master)
[![Code Climate](https://img.shields.io/codeclimate/github/raphink/narcissus.svg?maxAge=2592000)](https://codeclimate.com/github/raphink/narcissus/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/raphink/narcissus)](https://goreportcard.com/report/github.com/raphink/narcissus)

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

	// Modify UID
	user.UID = 42
  
	err = n.Write(user)
	if err != nil {
		log.Fatalf("Failed to save user: %v", err)
	}
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
	Name       string   `narcissus:".,value-from-label"`
	Password   string   `narcissus:"password"`
	GID        int      `narcissus:"gid"`
	Users      []string `narcissus:"user"`
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
	log.Printf("Users=%v", strings.Join(group.Users, ","))
}
```


## Tag values

The `narcissus` tag accepts multiple comma separated values. The first value in
the list is the relative path where the field is mapped in the Augeas tree.

Other possible (optional) values are:

* `value-from-label`: get field value from the node label instead of its value;
* `seq` (slice field only): will treat field as a seq entry in the Augeas tree;
* `key-from-value` (map field only): get the key from the node label instead
   of its value;
* `purge` (map field only): purge all unknown keys in the map.


## Fields

Described structures have special fields to specify parameters for Augeas.

* `augeasPath` (mandatory): the path to the structure in the Augeas tree;
* `augeasFile` (optional): let Augeas load only this file. If `augeasLens` is
  not specified, Augeas will use the default lens for the file if available;
* `augeasLens` (optional): if `augeasFile` is set (ignored otherwise),
  specifies which lens to use to parse the file. This is required when parsing
  a file at a non-standard location.


Each of these fields can be specified in one of two ways:

* by using the `default` tag with a default value for the field, e.g.

```go
type group struct {
	augeasPath string `default:"/files/etc/group/root"`
}
```

* by specifying a value for the instance in the structure field, e.g.

```go
myGroup := group {
    augeasPath: "/files/etc/group/docker",
}
```
