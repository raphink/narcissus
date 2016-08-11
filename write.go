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
	refType := ref.Type()
	for i := 0; i < refType.NumField(); i++ {
		if refType.Field(i).Name == "augeasPath" {
			// Ignore the special `augeasPath` field
			continue
		}
		err := n.writeField(ref.Field(i), refType.Field(i), path)
		if err != nil {
			return fmt.Errorf("failed to write field %s: %v", refType.Field(i).Name, err)
		}
	}

	return nil
}

func (n *Narcissus) writeField(field reflect.Value, fieldType reflect.StructField, path string) error {
	return fmt.Errorf("Write() is not implemented yet")
}
