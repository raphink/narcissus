package narcissus

// Fstab maps a /etc/fstab file
type Fstab struct {
	Entries []FstabEntry `type:"seq"`
}

// FstabEntry maps an Fstab entry
type FstabEntry struct {
	Spec    string `path:"spec"`
	File    string `path:"file"`
	Vfstype string `path:"vfstype"`
	Opt     []struct {
		Key   string `path:"."`
		Value string `path:"value"`
	} `path:"opt"`
	Dump   int `path:"dump"`
	Passno int `path:"passno"`
}
