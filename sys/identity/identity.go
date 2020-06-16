package identity

import (
	"log"
	"net"
)

var (
	Identity string
)

type IdentitySection struct {
	Specify string `yaml:"specify"`
	Shell   string `yaml:"shell"`
}

func MyIp4List() ([]string, error) {
	ips := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}

func Init(identity IdentitySection) {
	if identity.Specify != "" {
		Identity = identity.Specify
		return
	}
	ips, err := MyIp4List()
	if err != nil {
		log.Fatalln("cannot get identity: ", err)
	}
	if len(ips) == 0 {
		log.Fatalln("cannot get identity, no global unicast ip can found ")
	}
	Identity = ips[0]
	return
}
