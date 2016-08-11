package narcissus

// Passwd maps /etc/passwd
type Passwd struct {
	Users map[string]PasswdUser `path:"./*" key:"label"`
}

// PasswdUser maps a Passwd user
type PasswdUser struct {
	Account  string `path:"." value-from:"label"`
	Password string `path:"password"`
	Uid      int    `path:"uid"`
	Gid      int    `path:"gid"`
	Name     string `path:"name"`
	Home     string `path:"home"`
	Shell    string `path:"shell"`
}
