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
	refType := ref.Type()
	if pType, ok := refType.FieldByName("augeasPath"); ok {
		p := ref.FieldByName("augeasPath")
		if pp := p.String(); pp == "" {
			if defaultP := pType.Tag.Get("default"); defaultP != "" {
				return defaultP, nil
			} else {
				return "", fmt.Errorf("no augeasPath value and no default")
			}
		} else {
			return pp, nil
		}
	}
	return "", fmt.Errorf("no augeasPath field")
}
