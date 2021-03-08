package main

import (
	"net"
	"sync"
	"time"
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
}

type Configuration struct {
	contactTTL time.Duration
	freqhz     float64
}

type Block struct {
	samples [][2]float64
}
