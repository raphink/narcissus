package narcissus

import (
	"fmt"
	"reflect"
)

func structRef(val interface{}) (ref reflect.Value, err error) {
	ptr := reflect.ValueOf(val)
	if ptr.Kind() != reflect.Ptr {
		err = fmt.Errorf("not a ptr")
		return
	}
	ref = ptr.Elem()

	if ref.Kind() != reflect.Struct {
		err = fmt.Errorf("not a struct")
		return
	}

	return
}

func getPath(ref reflect.Value) (string, error) {
	return getField("augeasPath", ref)
}

func getFile(ref reflect.Value) (string, error) {
	return getField("augeasFile", ref)
}

func getLens(ref reflect.Value) (string, error) {
	return getField("augeasLens", ref)
}

func getField(name string, ref reflect.Value) (string, error) {
	refType := ref.Type()
	if pType, ok := refType.FieldByName(name); ok {
		p := ref.FieldByName(name)
		pp := p.String()
		if pp == "" {
			if defaultP := pType.Tag.Get("default"); defaultP != "" {
				return defaultP, nil
			}
			return "", fmt.Errorf("no %s value and no default", name)
		}
		return pp, nil
	}
	return "", fmt.Errorf("no %s field", name)
}
