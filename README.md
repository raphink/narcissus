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

func main() {
	aug, err := augeas.New("/", "", augeas.None)
	if err != nil {
		log.Fatal("Failed to create Augeas handler")
	}

	user := &narcissus.PasswdUser{}
	err = narcissus.Parse(aug, user, "/files/etc/passwd/raphink")
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	log.Printf("Uid=%v", user.Uid)
}
```



## Mapping your own structures


```go
import (
	"log"

	"honnef.co/go/augeas"
	"github.com/raphink/narcissus"
)

type Group struct {
	Name     string `path:"." value-from:"label"`
	Password string `path:"password"`
	Gid      int    `path:"gid"`
}


func main() {
	aug, err := augeas.New("/", "", augeas.None)
	if err != nil {
		log.Fatal("Failed to create Augeas handler")
	}

	group := &Group{}
	err = narcissus.Parse(aug, user, "/files/etc/group/docker")
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	log.Printf("Uid=%v", group.Gid)
}
```
