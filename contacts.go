package main

import (
	"fmt"
	"strconv"
	"time"
)

func (c Contact) String() string {
	return fmt.Sprintf("#%d %s %s (%s:%d)", c.index, c.Hostname, c.ID, c.Addr.IP, c.Addr.Port)
}

func (c *Contacts) check() {
	for uid, value := range c.contacts {
		if value.last.Add(value.ttl).Before(time.Now()) {
			c.mutex.Lock()
			delete(c.contacts, uid)
			fmt.Println("Contact expired:", value)
			c.mutex.Unlock()
		}
	}
}

func (c *Contacts) add(contact *Contact) {
	c.mutex.Lock()
	_, found := c.contacts[contact.ID]
	if !found {
		//fmt.Printf("add %p\n", contact)
		c.contacts[contact.ID] = contact
	} else {
		c.contacts[contact.ID].last = time.Now()
	}
	c.mutex.Unlock()
}

func (c *Contacts) Len() int {
	return len(c.contacts)
}

func (c *Contact) Stop() {
	c.bs.drain = true
	c.listening = false
}

func (contacts *Contacts) find(args []string) <-chan *Contact {
	ret := make(chan *Contact)
	go func() {
		if len(args) == 0 {
			for _, contact := range contacts.contacts {
				ret <- contact
			}
		}
		for _, c := range args {
			if c[0] == '#' {
				// by index
				index, err := strconv.Atoi(c[1:])
				if err != nil {
					return
				}
				for _, v := range contacts.contacts {
					if v.index == index {
						ret <- v
						break
					}
				}
			} else {
				// try the uid
				val, found := contacts.contacts[c]
				if found {
					ret <- val
				} else {
					// try the hostname
					for _, v := range contacts.contacts {
						if v.Hostname == c {
							ret <- v
							break
						}
					}
				}
			}
		}
		close(ret)
	}()
	return ret
}
