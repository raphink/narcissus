package narcissus

import (
	"fmt"
	"reflect"
	"strconv"

	"honnef.co/go/augeas"
)

func Parse(aug augeas.Augeas, val interface{}, path string) error {
	ptr := reflect.ValueOf(val)
	if ptr.Kind() != reflect.Ptr {
		return fmt.Errorf("not a ptr")
	}

	ref := ptr.Elem()
	if ref.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct ptr")
	}

	refType := ref.Type()
	for i := 0; i < refType.NumField(); i++ {
		value, err := getField(aug, refType.Field(i), path)
		if err != nil {
			return err
		}
		setField(ref.Field(i), refType.Field(i), value)
	}

	return nil
}

func getField(aug augeas.Augeas, field reflect.StructField, path string) (value string, err error) {
	fieldPath := fmt.Sprintf("%s/%s", path, field.Tag.Get("path"))
	value, err = aug.Get(fieldPath)
	return
}

func setField(field reflect.Value, fieldType reflect.StructField, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		bvalue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(bvalue)
	default:
		return fmt.Errorf("unsupported type %s", field.Kind())
	}
	return nil
}
