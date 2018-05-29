package narcissus

// Hosts maps a /etc/hosts file
type Hosts struct {
	augeasPath string `default:"/files/etc/hosts"`
	Comments   []struct {
		Comment string `narcissus:"."`
	} `narcissus:"#comment"`
	Hosts []Host `narcissus:"seq"`
}

// Host maps an Hosts entry
type Host struct {
	augeasPath string
	IPAddress  string   `narcissus:"ipaddr"`
	Canonical  string   `narcissus:"canonical"`
	Aliases    []string `narcissus:"alias"`
	Comment    string   `narcissus:"#comment"`
}

// NewHosts returns a new Hosts structure
func (n *Narcissus) NewHosts() (h *Hosts, err error) {
	h = &Hosts{}
	err = n.Parse(h)
	return
}
