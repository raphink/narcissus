package narcissus

import (
	"fmt"
	"reflect"
)

func (n *Narcissus) Autoload(ref reflect.Value) (err error) {
	if file, _ := getFile(ref); file != "" {
		if lens, _ := getLens(ref); lens != "" {
			err = n.Augeas.Transform(lens, file, false)
			if err != nil {
				return fmt.Errorf("failed to set up Augeas transformation for file %s using lens %s: %v", file, lens, err)
			}
		} else {
			err = n.Augeas.LoadFile(file)
			if err != nil {
				return fmt.Errorf("failed to load file %s in Augeas: %v", file, err)
			}
		}
		err = n.Augeas.Load()
		if err != nil {
			return fmt.Errorf("failed to load Augeas tree: %v", err)
		}

		errpath := "/augeas/files" + file + "/error/message"
		if msg, _ := n.Augeas.Get(errpath); msg != "" {
			return fmt.Errorf("failed to load file %s in Augeas: %v", msg)
		}
	}

	return
}
