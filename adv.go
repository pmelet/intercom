package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const broadcast_addr = "255.255.255.255"

func advertise(wg *sync.WaitGroup, ctx *shellContext, lstAddr net.TCPAddr) (err error) {
	defer wg.Done()
	defer ctx.Set("advertise", false)
	ctx.Set("advertise", true)

	advAddr, err := net.ResolveUDPAddr("udp", broadcast_addr+":4242")
	if err != nil {
		return
	}
	conn, err := net.DialUDP("udp", nil, advAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	hn, err := os.Hostname()
	if err != nil {
		return
	}
	m := Manifest{lstAddr, hn, ctx.ID}
	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	for {
		_, err2 := fmt.Fprintf(conn, string(b))
		if err != nil {
			return err2
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
