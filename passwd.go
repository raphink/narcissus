package narcissus

// Passwd maps /etc/passwd
type Passwd struct {
	Users map[string]PasswdUser `path:"*" key:"label"`
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

// NewPasswd returns a new Passwd structure
func (n *Narcissus) NewPasswd() (p *Passwd, err error) {
	p = &Passwd{}
	err = n.Parse(p, "/files/etc/passwd")
	return
}

// NewPasswdUser returns a new PasswdUser structure for given user
func (n *Narcissus) NewPasswdUser(user string) (p *PasswdUser, err error) {
	p = &PasswdUser{}
	err = n.Parse(p, "/files/etc/passwd/"+user)
	return
}
