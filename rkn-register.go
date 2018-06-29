package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
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
func (reg *Register) Load(r io.Reader) {
	reg.ip = bytetree.NewTree()

	s := bufio.NewScanner(r)
	s.Scan() //skip first line
	for s.Scan() {
		ips, err := reg.parseLine(s.Text())
		if err != nil {
			log.Println(err)
			continue
		}

		for _, ip := range ips {
			reg.ip.GrowLeaf(ip, true)
		}
	}
}

//parseLine from csv file, expected format: '1.173.24.102 | 61.227.155.222;;;...'
func (reg *Register) parseLine(line string) (ips []net.IP, err error) {
	n := strings.Index(line, ";")
	if n < 1 {
		return ips, errors.New("unexpected line, separator ';' not found")
	}

	for _, sip := range strings.Split(line[:n], "|") {
		ip := net.ParseIP(strings.TrimSpace(sip))
		if ip.IsUnspecified() {
			log.Println("invalid ip address", sip)
			continue
		}

		ips = append(ips, ip)
	}
	return
}

func (reg *Register) LookupIP(ip net.IP) bool {
	_, ok := reg.ip.LookupLeaf(ip)
	return ok
}
