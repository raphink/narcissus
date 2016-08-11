package narcissus

import (
	"fmt"
	"reflect"
)

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
