package narcissus

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Parse parses a structure pointer and feeds its fields with Augeas data
func (n *Narcissus) Parse(val interface{}) error {
	ptr := reflect.ValueOf(val)
	if ptr.Kind() != reflect.Ptr {
		return fmt.Errorf("not a ptr")
	}
	ref := ptr.Elem()

	if ref.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct")
	}

	refType := ref.Type()
	var path string
	if pType, ok := refType.FieldByName("augeasPath"); ok {
		p := ref.FieldByName("augeasPath")
		if pp := p.String(); pp == "" {
			if defaultP := pType.Tag.Get("default"); defaultP != "" {
				path = defaultP
			} else {
				return fmt.Errorf("no augeasPath value and no default")
			}
		} else {
			path = pp
		}
	} else {
		return fmt.Errorf("no augeasPath field")
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
		if err == ErrNodeNotFound {
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
	fieldPath := fmt.Sprintf("%s/%s", path, fieldType.Tag.Get("path"))
	if field.Kind() == reflect.Slice {
		return n.getSliceField(field, fieldType, path, fieldPath)
	} else if field.Kind() == reflect.Map {
		return n.getMapField(field, fieldType, fieldPath)
	} else {
		return n.getStringField(fieldType.Type, fieldPath, fieldType.Tag)
	}
	return nil, nil
}

func (n *Narcissus) getStringField(fieldType reflect.Type, fieldPath string, tag reflect.StructTag) (string, error) {
	aug := n.Augeas
	var value string
	var err error
	if tag.Get("value-from") == "label" {
		value, err = aug.Label(fieldPath)
	} else {
		value, err = aug.Get(fieldPath)
	}
	if err != nil && strings.Contains(err.Error(), "No matching node") {
		return "", ErrNodeNotFound
	}
	return value, err
}

func (n *Narcissus) getSliceField(field reflect.Value, fieldType reflect.StructField, path, fieldPath string) (interface{}, error) {
	aug := n.Augeas

	if fieldType.Tag.Get("type") == "seq" {
		fieldPath = fmt.Sprintf("%s/*[.=~regexp('[0-9]*')]", path)
	}
	matches, err := aug.Match(fieldPath)
	if err != nil {
		return nil, err
	}
	values := reflect.MakeSlice(field.Type(), len(matches), len(matches))
	for i, m := range matches {
		vType := field.Type().Elem()
		var value reflect.Value
		if vType.Kind() == reflect.Struct {
			vStruct := reflect.New(vType)
			err = n.parseStruct(vStruct.Elem(), m)
			if err != nil {
				return nil, fmt.Errorf("failed to get slice element: %v", err)
			}
			value = vStruct.Elem()
		} else {
			val, err := n.getStringField(vType, m, "")
			if err != nil {
				return nil, fmt.Errorf("failed to get slice element: %v", err)
			}
			value = reflect.ValueOf(val)
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
		if fieldType.Tag.Get("key") == "label" {
			label, err = aug.Label(m)
		} else {
			label, err = aug.Get(m)
		}
		if err != nil {
			return nil, err
		}
		vType := field.Type().Elem()
		vStruct := reflect.New(vType)
		err = n.parseStruct(vStruct.Elem(), m)
		if err != nil {
			return nil, err
		}
		values.SetMapIndex(reflect.ValueOf(label), vStruct.Elem())
	}
	return values.Interface(), nil
}

func setField(field reflect.Value, fieldType reflect.StructField, value interface{}) error {
	v := reflect.ValueOf(value)
	if field.Kind() == v.Kind() {
		field.Set(v)
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value.(string))
	case reflect.Bool:
		bvalue, err := strconv.ParseBool(value.(string))
		if err != nil {
			return err
		}
		field.SetBool(bvalue)
	case reflect.Int:
		intValue, err := strconv.ParseInt(value.(string), 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	default:
		return fmt.Errorf("unsupported type %s", field.Kind())
	}
	return nil
}
