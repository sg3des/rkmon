package main

import (
	"log"
	"net"
	"time"
)

type Subscriber struct {
	Email  string   `schema:"email,required"`
	IP     []net.IP `schema:"ip,required"`
	sentAt time.Time
}

func (subs *Subscriber) Key() []byte {
	return []byte(subs.Email)
}

//CheckSubscribers check all subscribers ip address and send email alert if any of them is blocked
func (api *API) CheckSubscribers() {
	for _, _subs := range api.subscribers.PickAllLeafs() {
		subs := _subs.(*Subscriber)

		blockedIP, blocked := api.CheckSubscriber(subs)
		if !blocked {
			continue
		}

		if subs.sentAt.Add(24 * time.Hour).After(time.Now()) {
			continue
		}

		if api.debug {
			log.Printf("send alert mail to %s, with IP-addresses: %v", subs.Email, blockedIP)
		}
		err := api.smtp.SendAlertTemplate(subs.Email, blockedIP)
		if err != nil {
			log.Printf("failed send mail to %s, by reason: %s", subs.Email, err)
			continue
		}
		subs.sentAt = time.Now()
	}
}

//CheckSubscriber check subscriber ip addresses and return lis of blocked ip addresses
func (api *API) CheckSubscriber(subs *Subscriber) (blockedIP []net.IP, blocked bool) {
	for _, ip := range subs.IP {
		if !api.reg.LookupIP(ip) {
			continue
		}

		blockedIP = append(blockedIP, ip)
	}

	blocked = len(blockedIP) > 0

	return
}

//AddSubscriber to checkable list
func (api *API) AddSubscriber(newsubs *Subscriber) {
	_subs, ok := api.subscribers.LookupLeaf(newsubs.Key())
	if ok {
		subs := _subs.(Subscriber)
		subs.IP = newsubs.IP
	} else {
		api.subscribers.GrowLeaf([]byte(newsubs.Email), newsubs)
	}
}

//RemoveSubscriber from checkable list
func (api *API) RemoveSubscriber(email string) {
	api.subscribers.CutLeaf([]byte(email))
}
