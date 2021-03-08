package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

func mainLoop(conf *Configuration, server bool, client bool) {
	id := uuid.NewString()
	contacts := Contacts{contacts: make(map[string]*Contact)}
	context := shellContext{contacts: &contacts, ID: id, keys: make(map[string]bool), conf: conf}

	wg := sync.WaitGroup{}
	nbRoutines := 1
	if server {
		nbRoutines += 2
	}
	if client {
		nbRoutines++
	}
	wg.Add(nbRoutines)
	if server {
		listener, lstAddr, err := getListenSocket()
		if err != nil {
			fmt.Println(err)
			return
		}
		// listen for clients
		go serve(&wg, &context, listener)
		// advertise our presence
		go advertise(&wg, &context, lstAddr)
	}
	if client {
		// listen for contacts on the LAN
		go listenAdv(&wg, &context)
	}
	// shell
	go shell(&wg, &context)
	wg.Wait()
}

/*
func callback(p *argparse.Parser, ns *argparse.Namespace, leftovers []string, err error) {
	if err != nil {
		switch err.(type) {
		case argparse.ShowHelpErr, argparse.ShowVersionErr:
			// For either ShowHelpErr or ShowVersionErr, the parser has already
			// displayed the necessary text to the user. So we end the program
			// by returning.
			return
		default:
			fmt.Println(err)
			p.ShowHelp()
		}
		return // Exit program
	}

	fmt.Println(ns.Get("server"))
	fmt.Println(ns.Get("client"))
	fmt.Println(leftovers)

	mainLoop(ns.Get("server").(string) == "true", ns.Get("client").(string) == "true")
}
*/

func main() {
	//p := argparse.NewParser("Intercom", callback).Version("0.0.0")
	//p.AddHelp().AddVersion() // Enable `--help` & `-h` to display usage text to the user.
	//server := argparse.NewFlag("s server", "server", "Emit sound, advertise presence")
	//client := argparse.NewFlag("c client", "client", "Detect presence, listen to sound")
	//p.AddOptions(server, client)
	//p.Parse(os.Args[1:]...)

	server := flag.Bool("s", false, "Server")
	client := flag.Bool("c", false, "Client")
	freqhz := flag.Float64("f", 220.0, "Frequency emitted")
	flag.Parse()

	conf := Configuration{contactTTL: time.Minute}
	conf.freqhz = *freqhz

	mainLoop(&conf, *server, *client)
}
