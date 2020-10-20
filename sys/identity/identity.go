package identity

import (
	"errors"
	"log"

	"net"
)

type Identity struct {
	IP    string `yaml:"ip"`
	Ident string `yaml:"ident"`
}

var config Identity

type IdentitySection struct {
	Specify string `yaml:"specify"`
}

type IPSection struct {
	Specify string `yaml:"specify"`
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

func getMyIP() (ip string, err error) {
	ips, err := MyIp4List()
	if err != nil {
		return
	}
	if len(ips) == 0 {
		err = errors.New("cannot get identity, no global unicast ip can found")
		return
	}
	ip = ips[0]
	return
}

func getIdent(identity IdentitySection) (string, error) {
	if identity.Specify != "" {
		return identity.Specify, nil
	}
	myip, err := getMyIP()
	if err != nil {
		return "", nil
	}
	return myip, nil
}

func getIP(ip IPSection) (string, error) {
	if ip.Specify != "" {
		return ip.Specify, nil
	}
	myip, err := getMyIP()
	if err != nil {
		return "", nil
	}
	return myip, nil
}

func Init(identity IdentitySection, ip IPSection) {
	ident, err := getIdent(identity)
	if err != nil {
		log.Fatalf("init identity failed, %v", err)
	}
	myip, err := getIP(ip)
	if err != nil {
		log.Fatalf("init ip failed, %v", err)
	}
	config.Ident = ident
	config.IP = myip
	return
}

func GetIP() string {
	return config.IP
}
func GetIdent() string {
	return config.Ident
}
