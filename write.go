package narcissus

import (
	"fmt"
	"reflect"
	"strings"
)

// Write writes a structure pointer to the Augeas tree
func (n *Narcissus) Write(val interface{}) error {
	ref, err := structRef(val)
	if err != nil {
		return fmt.Errorf("invalid interface: %v", err)
	}

	err = n.Autoload(ref)
	if err != nil {
		return fmt.Errorf("failed to autoload file: %v", err)
	}

	path, err := getPath(ref)
	if err != nil {
		return fmt.Errorf("undefined path: %v", err)
	}

	err = n.writeStruct(ref, path)
	if err != nil {
		return fmt.Errorf("failed to write interface to tree: %v", err)
	}

	err = n.Augeas.Save()
	if err != nil {
		return fmt.Errorf("failed to save Augeas tree: %v", err)
	}

	return nil
}

func (n *Narcissus) writeStruct(ref reflect.Value, path string) error {
	refType := ref.Type()
	for i := 0; i < refType.NumField(); i++ {
		if strings.HasPrefix(refType.Field(i).Name, "augeas") {
			// Ignore the special `augeas*` fields
			continue
		}
		err := n.writeField(ref.Field(i), refType.Field(i), path)
		if err != nil {
			return fmt.Errorf("failed to write field %s to path %s: %v", refType.Field(i).Name, path, err)
		}
	}

	return nil
}

func (n *Narcissus) writeField(field reflect.Value, fieldType reflect.StructField, path string) error {
	fieldPath := fmt.Sprintf("%s/%s", path, parseTag(fieldType.Tag).path)
	if field.Kind() == reflect.Slice {
		return n.writeSliceField(field, fieldType, path, fieldPath)
	} else if field.Kind() == reflect.Map {
		return n.writeMapField(field, fieldType, path, fieldPath)
	}
	return n.writeSimpleField(field, fieldPath, fieldType.Tag)
}

func (n *Narcissus) writeSimpleField(field reflect.Value, fieldPath string, tag reflect.StructTag) error {
	aug := n.Augeas
	// There might be a better way to convert, but that does it
	value := fmt.Sprintf("%v", field.Interface())

	if parseTag(tag).valueFromLabel {
		return nil
	}

	if value == "" && parseTag(tag).omitEmpty {
		return nil
	}

	err := aug.Set(fieldPath, value)
	return err
}

func (n *Narcissus) writeSliceField(field reflect.Value, fieldType reflect.StructField, path, fieldPath string) error {
	for i := 0; i < field.Len(); i++ {
		value := field.Index(i)
		var p string
		if parseTag(fieldType.Tag).seq {
			p = fmt.Sprintf("%s/%v", path, i+1)
		} else {
			p = fmt.Sprintf("%s[%v]", fieldPath, i+1)
		}
		// Create base node
		err := n.Augeas.Clear(p)
		if err != nil {
			return fmt.Errorf("failed to create base path for slice: %v", err)
		}
		if value.Kind() == reflect.Struct {
			err := n.writeStruct(value, p)
			if err != nil {
				return fmt.Errorf("failed to write slice struct value: %v", err)
			}
		} else {
			err := n.writeSimpleField(value, p, fieldType.Tag)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %v", err)
			}
		}
	}

	return nil
}

func (n *Narcissus) writeMapField(field reflect.Value, fieldType reflect.StructField, path, fieldPath string) (err error) {
	keys := field.MapKeys()
	var purgeConditions []string
	tag := fieldType.Tag
	for _, k := range keys {
		value := field.MapIndex(k)
		var p string
		if strings.HasSuffix(fieldPath, "/*") {
			// TrimSuffix? ouch!
			p = fmt.Sprintf("%s/%s", strings.TrimSuffix(fieldPath, "/*"), k)
			purgeConditions = append(purgeConditions, fmt.Sprintf("label() != '%s'", k))
		} else {
			p = fmt.Sprintf("%s[.='%s']", fieldPath, k)
			purgeConditions = append(purgeConditions, fmt.Sprintf(". != '%s'", k))
		}
		if value.Kind() == reflect.Struct {
			if !parseTag(tag).keyFromLabel {
				err = n.writeSimpleField(k, p, tag)
				if err != nil {
					return fmt.Errorf("failed to write map key: %v", err)
				}
			}
			err = n.writeStruct(value, p)
			if err != nil {
				return fmt.Errorf("failed to write map struct value: %v", err)
			}
		} else {
			err := n.writeSimpleField(value, p, tag)
			if err != nil {
				return fmt.Errorf("failed to write map value: %v", err)
			}
		}
	}

	// Purge absent keys
	if parseTag(tag).purge {
		purgePath := fieldPath + "[" + strings.Join(purgeConditions, " and ") + "]"
		n.Augeas.Remove(purgePath)
	}

	return nil
}
