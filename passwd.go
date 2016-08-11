package narcissus

// Passwd maps /etc/passwd
type Passwd struct {
	augeasPath string                `default:"/files/etc/passwd"`
	Users      map[string]PasswdUser `path:"*" key:"label"`
}

// PasswdUser maps a Passwd user
type PasswdUser struct {
	augeasPath string
	Account    string `path:"." value-from:"label"`
	Password   string `path:"password"`
	UID        int    `path:"uid"`
	GID        int    `path:"gid"`
	Name       string `path:"name"`
	Home       string `path:"home"`
	Shell      string `path:"shell"`
}

// NewPasswd returns a new Passwd structure
func (n *Narcissus) NewPasswd() (p *Passwd, err error) {
	p = &Passwd{}
	err = n.Parse(p)
	return
}

// NewPasswdUser returns a new PasswdUser structure for given user
func (n *Narcissus) NewPasswdUser(user string) (p *PasswdUser, err error) {
	p = &PasswdUser{
		augeasPath: "/files/etc/passwd/" + user,
	}
	err = n.Parse(p)
	return
}
