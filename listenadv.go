package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

const maxBufferSize = 1024

func listenAdv(wg *sync.WaitGroup, ctx *shellContext) error {
	defer wg.Done()
	defer ctx.Set("contacts", false)
	ctx.Set("contacts", true)

	var err error
	addr, err := net.ResolveUDPAddr("udp", ":4242")
	if err != nil {
		fmt.Println(err)
		return err
	}
	pc, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	buffer := make([]byte, maxBufferSize)
	for {
		n, udpAddr, err := pc.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return err
		}
		m := new(Manifest)
		err = json.Unmarshal(buffer[0:n], m)

		lstAddr := net.TCPAddr{IP: udpAddr.IP, Port: m.Addr.Port, Zone: m.Addr.Zone}

		contact := Contact{Addr: lstAddr, Hostname: m.Hostname, ID: m.ID, last: time.Now(), ttl: ctx.conf.contactTTL}
		ctx.contacts.check()
		ctx.contacts.add(&contact)
	}
}
