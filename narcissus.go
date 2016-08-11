package narcissus

import (
	"fmt"
	"log"
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
		value, err := getField(aug, ref.Field(i), refType.Field(i), path)
		if err != nil {
			return err
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

func getField(aug augeas.Augeas, field reflect.Value, fieldType reflect.StructField, path string) (interface{}, error) {
	fieldPath := fmt.Sprintf("%s/%s", path, fieldType.Tag.Get("path"))
	if field.Kind() == reflect.Slice {
		var values []interface{}
		matches, err := aug.Match(fieldPath)
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			v, err := aug.Get(m)
			if err != nil {
				return nil, err
			}
			values = append(values, v)
		}
		return nil, nil
	} else {
		log.Printf("Getting %s", fieldPath)
		var value string
		var err error
		if fieldType.Tag.Get("value-from") == "label" {
			value, err = aug.Label(fieldPath)
		} else {
			value, err = aug.Get(fieldPath)
		}
		return value, err
	}
	return nil, nil
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
