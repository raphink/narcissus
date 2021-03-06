package narcissus

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Parse parses a structure pointer and feeds its fields with Augeas data
func (n *Narcissus) Parse(val interface{}) error {
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

	return n.parseStruct(ref, path)
}

func (n *Narcissus) parseStruct(ref reflect.Value, path string) error {
	refType := ref.Type()
	for i := 0; i < refType.NumField(); i++ {
		if refType.Field(i).Name == "augeasPath" {
			// Ignore the special `augeasPath` field
			continue
		}
		value, err := n.getField(ref.Field(i), refType.Field(i), path)
		// TODO: implement omitempty
		if err == errNodeNotFound {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to retrieve field %s: %v", refType.Field(i).Name, err)
		}
		if value == "" {
			// for now
			continue
		}
		err = setField(ref.Field(i), refType.Field(i), value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Narcissus) getField(field reflect.Value, fieldType reflect.StructField, path string) (interface{}, error) {
	fieldPath := fmt.Sprintf("%s/%s", path, parseTag(fieldType.Tag).path)
	if field.Kind() == reflect.Slice {
		return n.getSliceField(field, fieldType, path, fieldPath)
	} else if field.Kind() == reflect.Map {
		return n.getMapField(field, fieldType, fieldPath)
	}
	return n.getSimpleField(fieldType.Type, fieldPath, fieldType.Tag)
}

func (n *Narcissus) getSimpleField(fieldType reflect.Type, fieldPath string, tag reflect.StructTag) (interface{}, error) {
	aug := n.Augeas
	var value string
	var err error
	if parseTag(tag).valueFromLabel {
		value, err = aug.Label(fieldPath)
	} else {
		value, err = aug.Get(fieldPath)
	}
	if err != nil && strings.Contains(err.Error(), "No matching node") {
		return nil, errNodeNotFound
	}

	switch fieldType.Kind() {
	case reflect.String:
		return value, nil
	case reflect.Bool:
		bvalue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to bool: %v", value, err)
		}
		return bvalue, nil
	case reflect.Int:
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to int: %v", value, err)
		}
		return int(intValue), nil
	}
	return nil, fmt.Errorf("unsupported type %s", fieldType.Kind())
}

func (n *Narcissus) getSliceField(field reflect.Value, fieldType reflect.StructField, path, fieldPath string) (interface{}, error) {
	aug := n.Augeas

	if parseTag(fieldType.Tag).seq {
		fieldPath = fmt.Sprintf("%s/*[label()=~regexp('[0-9]*')]", path)
	}
	matches, err := aug.Match(fieldPath)
	if err != nil {
		return nil, err
	}
	values := reflect.MakeSlice(field.Type(), len(matches), len(matches))
	for i, m := range matches {
		value, err := n.getStructOrSimpleField(field, m)
		if err != nil {
			return nil, fmt.Errorf("Failed to get slice element: %v", err)
		}
		values.Index(i).Set(value)
	}
	return values.Interface(), nil
}

func (n *Narcissus) getMapField(field reflect.Value, fieldType reflect.StructField, fieldPath string) (interface{}, error) {
	aug := n.Augeas
	values := reflect.MakeMap(field.Type())
	matches, err := aug.Match(fieldPath)
	if err != nil {
		return nil, err
	}
	for _, m := range matches {
		var label string
		if parseTag(fieldType.Tag).keyFromLabel {
			label, err = aug.Label(m)
		} else {
			label, err = aug.Get(m)
		}
		if err != nil {
			return nil, err
		}
		value, err := n.getStructOrSimpleField(field, m)
		if err != nil {
			return nil, fmt.Errorf("Failed to get map element: %v", err)
		}
		values.SetMapIndex(reflect.ValueOf(label), value)
	}
	return values.Interface(), nil
}

func (n *Narcissus) getStructOrSimpleField(field reflect.Value, path string) (value reflect.Value, err error) {
	vType := field.Type().Elem()
	if vType.Kind() == reflect.Struct {
		vStruct := reflect.New(vType)
		err = n.parseStruct(vStruct.Elem(), path)
		if err != nil {
			err = fmt.Errorf("failed to get element: %v", err)
			return
		}
		return vStruct.Elem(), nil
	}
	val, err := n.getSimpleField(vType, path, "")
	if err != nil {
		err = fmt.Errorf("failed to get element: %v", err)
		return
	}
	return reflect.ValueOf(val), nil
}

func setField(field reflect.Value, fieldType reflect.StructField, value interface{}) error {
	if !field.CanSet() {
		return fmt.Errorf("cannot set field %s", fieldType.Name)
	}

	v := reflect.ValueOf(value)
	if field.Kind() == v.Kind() {
		field.Set(v)
		return nil
	}

	return fmt.Errorf("Wrong value type (%v for %v)", v.Kind(), field.Kind())
}
