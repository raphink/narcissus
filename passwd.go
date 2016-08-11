package narcissus

// Passwd maps /etc/passwd
type Passwd struct {
	Users []PasswdUser
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
