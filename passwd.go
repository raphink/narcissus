package narcissus

// Passwd maps /etc/passwd
type Passwd struct {
	augeasPath string                `default:"/files/etc/passwd"`
	Users      map[string]PasswdUser `narcissus:"*,key-from-label"`
}

// PasswdUser maps a Passwd user
type PasswdUser struct {
	augeasPath string
	Account    string `narcissus:".,value-from-label"`
	Password   string `narcissus:"password"`
	UID        int    `narcissus:"uid"`
	GID        int    `narcissus:"gid"`
	Name       string `narcissus:"name"`
	Home       string `narcissus:"home"`
	Shell      string `narcissus:"shell"`
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
