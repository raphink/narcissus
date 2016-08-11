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

func getField(aug augeas.Augeas, field reflect.Value, fieldType reflect.StructField, path string) (value string, err error) {
	if field.Kind() == reflect.Slice {
		log.Println("Unsupported type slice")
		return
	} else {
		fieldPath := fmt.Sprintf("%s/%s", path, fieldType.Tag.Get("path"))
		log.Printf("Getting %s", fieldPath)
		value, err = aug.Get(fieldPath)
	}
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
	case reflect.Int:
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	default:
		return fmt.Errorf("unsupported type %s", field.Kind())
	}
	return nil
}
