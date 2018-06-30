package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/sg3des/bytetree"
)

type Register struct {
	ip *bytetree.Tree
}

//NewRegister initialize storage for ipaddresses
func NewRegister() *Register {
	return &Register{ip: bytetree.NewTree()}
}

//Load new data to Register
func (reg *Register) Load(f *os.File) {
	reg.ip = bytetree.NewTree()

	r := bufio.NewReader(f)
	r.ReadLine()
	for i := 0; ; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		// log.Println(string(line))
		// time.Sleep(500 * time.Millisecond)

		ips, err := reg.parseLine(string(line))
		if err != nil {
			log.Println(i, len(line), err)
		}

		for _, ip := range ips {
			reg.ip.GrowLeaf(ip, true)
		}
	}
}

//parseLine from csv file, expected format: '1.173.24.102 | 61.227.155.222;;;...'
func (reg *Register) parseLine(line string) (ips []net.IP, err error) {
	n := strings.Index(line, ";")
	if n < 0 {
		return ips, errors.New("unexpected line, separator ';' not found " + line)
	}

	for _, s := range strings.Split(line[:n], "|") {
		s = strings.TrimSpace(s)

		//expand CIDR to multiple IP addresses
		if strings.Contains(s, "/") {
			cidrips, err := reg.expandCIDR(s)
			if err != nil {
				log.Println("failed expand CIDR:", err, "line:", line)
				continue
			}
			ips = append(ips, cidrips...)
			continue
		}

		//parse string to net.IP
		ip := net.ParseIP(s)
		if ip.IsUnspecified() {
			log.Println("invalid ip address", s)
			continue
		}

		ips = append(ips, ip)
	}
	return
}

func (reg *Register) expandCIDR(cidr string) (ips []net.IP, err error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); reg.ipInclude(ip) {
		ips = append(ips, ip)
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func (reg *Register) ipInclude(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func (reg *Register) LookupIP(ip net.IP) bool {
	_, ok := reg.ip.LookupLeaf(ip)
	return ok
}

func (reg *Register) TotalIP() int {
	return len(reg.ip.PickAllLeafs())
}
