Narcissus
=========

This go package aims to provide reflection for the Augeas library.

## Example

```go
import (
	"log"

	"honnef.co/go/augeas"
	"github.com/raphink/narcissus"
)

type PasswdUser struct {
	Account  string `path:"." value-from:"label"`
	Password string `path:"password"`
	Uid      int    `path:"uid"`
	Gid      int    `path:"gid"`
	Name     string `path:"name"`
	Home     string `path:"home"`
	Shell    string `path:"shell"`
}


func main() {
	aug, err := augeas.New("/", "", augeas.None)
	if err != nil {
		log.Fatal("Failed to create Augeas handler")
	}

	user := &PasswdUser{}
	err = narcissus.Parse(aug, user, "/files/etc/passwd/raphink")
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

  log.Printf("Uid=%s", user.Uid)
}
```

