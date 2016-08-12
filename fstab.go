package narcissus

// Fstab maps a /etc/fstab file
type Fstab struct {
	augeasPath string `default:"/files/etc/fstab"`
	Comments   []struct {
		Comment string `path:"."`
	} `path:"#comment"`
	Entries []FstabEntry `type:"seq"`
}

// FstabEntry maps an Fstab entry
type FstabEntry struct {
	augeasPath string
	Spec       string              `path:"spec"`
	File       string              `path:"file"`
	Vfstype    string              `path:"vfstype"`
	Opt        map[string]FstabOpt `path:"opt"`
	Dump       int                 `path:"dump"`
	Passno     int                 `path:"passno"`
}

// FstabOpt is an FstabEntry opt
type FstabOpt struct {
	Value string `path:"value"`
}

// NewFstab returns a new Fstab structure
func (n *Narcissus) NewFstab() (f *Fstab, err error) {
	f = &Fstab{}
	err = n.Parse(f)
	return
}
