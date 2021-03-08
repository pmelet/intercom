package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/faiface/beep"
)

type actionCommand func(ctx *shellContext, args []string) error

type action struct {
	cmd actionCommand
	sub map[string]action
}

var commands = map[string]action{
	"status": {
		cmd: statusCmd,
	},
	"list": {
		cmd: listContacts,
	},
	"listen": {
		cmd: startListenContact,
	},
	"mute": {
		cmd: stopListenContact,
	},
	"test": {
		sub: map[string]action{
			"sound": {
				sub: map[string]action{
					"start": {cmd: playSound},
					"stop":  {cmd: muteSound},
				},
			},
		},
	},
}

func getKeys(m map[string]action) (keys []string) {
	keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func shell(wg *sync.WaitGroup, ctx *shellContext) {
	defer wg.Done()
	defer ctx.Set("shell", false)
	ctx.Set("shell", true)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		command = strings.Trim(command, " \n")
		c := regexp.MustCompile(" +").Split(command, -1)

		var p map[string]action = commands
		var execute actionCommand = nil
		var syntaxError = false
		var lastIndex int = 0
	Loop:
		for i, x := range c {
			if len(x) == 0 {
				// empty string
				continue
			}
			lastIndex = i
			entry, found := p[x]
			switch {
			case !found:
				syntaxError = true
				break Loop
			case entry.sub != nil:
				p = entry.sub
			case entry.cmd != nil:
				execute = entry.cmd
				break Loop
			default:
				syntaxError = true
				break Loop
			}
		}
		switch {
		case syntaxError:
			fmt.Println("*** error")
		case execute == nil:
			fmt.Println("Available commands:")
			for _, k := range getKeys(p) {
				fmt.Println("-", strings.Join(c[:lastIndex+1], " "), k)
			}
		default:
			execute(ctx, c[lastIndex+1:])
		}
	}
}

func statusCmd(ctx *shellContext, args []string) (err error) {
	fmt.Println(ctx.keys)
	if ctx.Get("advertise") {
		fmt.Println("Advertise")
	} else {
		fmt.Println("Do not advertise")
	}
	if ctx.Get("serve") {
		fmt.Println("Serve")
	} else {
		fmt.Println("Do not serve")
	}
	if ctx.Get("contacts") {
		fmt.Println("Looking for contacts")
	} else {
		fmt.Println("Not looking for contacts")
	}
	if ctx.Get("shell") {
		fmt.Println("Shell is running")
	} else {
		fmt.Println("Shell is not running (?!)")
	}
	return
}

var sound chan bool

func playSound(ctx *shellContext, args []string) (err error) {
	sound = make(chan bool)
	go play(sound)
	return
}

func muteSound(ctx *shellContext, args []string) (err error) {
	sound <- true
	return
}

func listContacts(ctx *shellContext, args []string) (err error) {
	for k, v := range ctx.contacts.contacts {
		l := " "
		if v.listening {
			l = "*"
		}
		fmt.Printf("%s %s %s %p\n", l, k, v, v)
	}
	return
}

func startListenContact(ctx *shellContext, args []string) (err error) {
	var contact *Contact
	contacts, err := ctx.contacts.find(args)
	if err != nil {
		return
	}
	for contact = range contacts {
		fmt.Printf("listen %p\n", contact)
		go listenContact(contact)
	}
	return
}

func stopListenContact(ctx *shellContext, args []string) (err error) {
	var contact *Contact
	contacts, err := ctx.contacts.find(args)
	if err != nil {
		return
	}
	for contact = range contacts {
		fmt.Printf("mute %p\n", contact)
		contact.Stop()
	}
	return
}

func listenContact(contact *Contact) error {
	conn, err := net.DialTCP("tcp4", nil, &contact.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	sr := beep.SampleRate(44100)
	contact.bs = BufferedStreamer{SampleRate: sr}
	contact.listening = true

	go playStream(&contact.bs, sr)

	for {
		var left float64
		var right float64
		err = binary.Read(conn, binary.BigEndian, &left)
		if err != nil {
			break
		}
		err = binary.Read(conn, binary.BigEndian, &right)
		if err != nil {
			break
		}
		contact.bs.Add(left, right)
	}
	return nil
}
