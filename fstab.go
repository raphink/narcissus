package narcissus

// Fstab maps a /etc/fstab file
type Fstab struct {
	Comments []struct {
		Comment string `path:"."`
	} `path:"#comment"`
	Entries []FstabEntry `type:"seq"`
}

// FstabEntry maps an Fstab entry
type FstabEntry struct {
	Spec    string `path:"spec"`
	File    string `path:"file"`
	Vfstype string `path:"vfstype"`
	Opt     map[string]struct {
		Value string `path:"value"`
	} `path:"opt"`
	Dump   int `path:"dump"`
	Passno int `path:"passno"`
}

// NewFstab returns a new Fstab structure
func (n *Narcissus) NewFstab() (f *Fstab, err error) {
	f = &Fstab{}
	err = n.Parse(f, "/files/etc/fstab")
	return
}
