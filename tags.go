package narcissus

import (
	"reflect"
	"strings"
)

const tagName = "narcissus"

type fieldTags struct {
	path           string
	purge          bool
	keyFromLabel   bool
	valueFromLabel bool
	seq            bool
}

func parseTag(tag reflect.StructTag) fieldTags {
	slice := strings.Split(tag.Get(tagName), ",")

	return fieldTags{
		path:           slice[0],
		purge:          sliceHasTag(slice, "purge"),
		keyFromLabel:   sliceHasTag(slice, "key-from-label"),
		valueFromLabel: sliceHasTag(slice, "value-from-label"),
		seq:            sliceHasTag(slice, "seq"),
	}
}

func sliceHasTag(slice []string, name string) bool {
	for _, t := range slice {
		if t == name {
			return true
		}
	}
	return false
}
