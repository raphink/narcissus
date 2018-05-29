package narcissus

// Fstab maps a /etc/fstab file
type Fstab struct {
	augeasPath string `default:"/files/etc/fstab"`
	Comments   []struct {
		Comment string `narcissus:"."`
	} `narcissus:"#comment"`
	Entries []FstabEntry `narcissus:"seq"`
}

// FstabEntry maps an Fstab entry
type FstabEntry struct {
	augeasPath string
	Spec       string              `narcissus:"spec"`
	File       string              `narcissus:"file"`
	Vfstype    string              `narcissus:"vfstype"`
	Opt        map[string]FstabOpt `narcissus:"opt"`
	Dump       int                 `narcissus:"dump"`
	Passno     int                 `narcissus:"passno"`
}

// FstabOpt is an FstabEntry opt
type FstabOpt struct {
	Value string `narcissus:"value"`
}

// NewFstab returns a new Fstab structure
func (n *Narcissus) NewFstab() (f *Fstab, err error) {
	f = &Fstab{}
	err = n.Parse(f)
	return
}
