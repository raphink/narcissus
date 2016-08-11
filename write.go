package narcissus

import (
	"fmt"
	"reflect"
)

// Write writes a structure pointer to the Augeas tree
func (n *Narcissus) Write(val interface{}) error {
	ref, err := structRef(val)
	if err != nil {
		return fmt.Errorf("invalid interface: %v", err)
	}

	path, err := getPath(ref)
	if err != nil {
		return fmt.Errorf("undefined path: %v", err)
	}

	return n.writeStruct(ref, path)
}

func (n *Narcissus) writeStruct(ref reflect.Value, path string) error {
	return fmt.Errorf("Write() is not implemented yet")
}
