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
	fieldPath := fmt.Sprintf("%s/%s", path, fieldType.Tag.Get("path"))
	if field.Kind() == reflect.Slice {
		return n.writeSliceField(field, fieldType, path, fieldPath)
	} else if field.Kind() == reflect.Map {
		return n.writeMapField(field, fieldType, fieldPath)
	}
	return n.writeSimpleField(fieldType.Type, fieldPath, fieldType.Tag)
}

func (n *Narcissus) writeSimpleField(fieldType reflect.Type, fieldPath string, tag reflect.StructTag) error {
	return nil
}

func (n *Narcissus) writeSliceField(field reflect.Value, fieldType reflect.StructField, path, fieldPath string) error {
	return nil
}

func (n *Narcissus) writeMapField(field reflect.Value, fieldType reflect.StructField, fieldPath string) error {
	return nil
}
