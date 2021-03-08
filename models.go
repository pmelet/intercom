package main

import (
	"net"
	"sync"
	"time"

	"github.com/faiface/beep"
)

type Manifest struct {
	Addr     net.TCPAddr
	Hostname string
	ID       string
}

type Contact struct {
	Addr      net.TCPAddr
	Hostname  string
	ID        string
	index     int
	last      time.Time
	ttl       time.Duration
	bs        BufferedStreamer
	listening bool
}

type Contacts struct {
	contacts map[string]*Contact
	mutex    sync.Mutex
}

type shellContext struct {
	contacts *Contacts
	ID       string
	mutex    sync.Mutex
	keys     map[string]bool
	conf     *Configuration
	mixer    *beep.Mixer
}

type Configuration struct {
	contactTTL time.Duration
	freqhz     float64
	port       int
}

type Block struct {
	samples [][2]float64
}
