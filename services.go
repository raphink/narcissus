package narcissus

// Services maps a /etc/services file
type Services struct {
	augeasPath string `default:"/files/etc/services"`
	Comments   []struct {
		Comment string `path:"."`
	} `path:"#comment"`
	Services []Service `path:"service-name"`
}

// Service maps a Services entry
type Service struct {
	augeasPath string
	Name       string `path:"."`
	Port       int    `path:"port"`
	Protocol   string `path:"protocol"`
	Comment    string `path:"#comment"`
}

// NewServices returns a new Services structure
func (n *Narcissus) NewServices() (s *Services, err error) {
	s = &Services{}
	err = n.Parse(s)
	return
}

// NewService returns a new Service structure
func (n *Narcissus) NewService(name string, protocol string) (s *Service, err error) {
	s = &Service{
		augeasPath: "/files/etc/services/service-name[.='" + name + "' and protocol='" + protocol + "']",
	}
	err = n.Parse(s)
	return
}
