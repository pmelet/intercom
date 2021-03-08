package main

import (
	"fmt"
	"time"
)

func (c Contact) String() string {
	return fmt.Sprintf("%s (%s:%d)", c.Hostname, c.Addr.IP, c.Addr.Port)
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
		fmt.Printf("add %p\n", contact)
		c.contacts[contact.ID] = contact
	} else {
		c.contacts[contact.ID].last = time.Now()
	}
	c.mutex.Unlock()
}

func (c *Contact) Stop() {
	c.bs.drain = true
	c.listening = false
}

func (contacts *Contacts) find(args []string) (<-chan *Contact, error) {
	ret := make(chan *Contact)
	go func() {
		if len(args) == 0 {
			for _, contact := range contacts.contacts {
				ret <- contact
			}
		}
		for _, c := range args {
			val, found := contacts.contacts[c]
			if found {
				ret <- val
			} else {
				for _, v := range contacts.contacts {
					if v.Hostname == c {
						ret <- v
						break
					}
				}
			}
		}
		close(ret)
	}()
	return ret, nil
}
