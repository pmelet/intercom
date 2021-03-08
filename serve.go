package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/faiface/beep"
)

func getListenSocket() (listener net.Listener, addr net.TCPAddr, err error) {
	lstAddr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		return
	}
	l, err := net.Listen("tcp4", lstAddr.String())
	fmt.Println(l.Addr())
	lstAddr, err = net.ResolveTCPAddr("tcp4", l.Addr().String())
	if err != nil {
		return
	}
	return l, *lstAddr, err
}

func serve(wg *sync.WaitGroup, ctx *shellContext, listener net.Listener) error {
	defer wg.Done()
	defer ctx.Set("serve", false)
	ctx.Set("serve", true)

	for {
		c, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return err
		}
		go handleConnection(ctx.conf, c)
	}
}

func handleConnection(conf *Configuration, c net.Conn) {
	defer c.Close()
	fmt.Println("Accepted connection")

	sr := beep.SampleRate(44100)
	w := &Wave{SampleRate: sr, freq: conf.freqhz}

	samples := make([][2]float64, 44100)
	for {
		w.Stream(samples)
		var buf bytes.Buffer
		err := binary.Write(&buf, binary.BigEndian, samples)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
		_, err = c.Write(buf.Bytes())
		if err != nil {
			fmt.Println("Connection closed", err)
			break
		}
		time.Sleep(time.Second * 3 / 4)
	}
}
