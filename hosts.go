package narcissus

// Hosts maps a /etc/hosts file
type Hosts struct {
	augeasPath string `default:"/files/etc/hosts"`
	Comments   []struct {
		Comment string `path:"."`
	} `path:"#comment"`
	Hosts []Host `type:"seq"`
}

// Host maps an Hosts entry
type Host struct {
	augeasPath string
	IPAddress  string   `path:"ipaddr"`
	Canonical  string   `path:"canonical"`
	Aliases    []string `path:"alias"`
	Comment    string   `path:"#comment"`
}

// NewHosts returns a new Hosts structure
func (n *Narcissus) NewHosts() (h *Hosts, err error) {
	h = &Hosts{}
	err = n.Parse(h)
	return
}
