package d7024e

import "sync"

type ContactChannel struct {
	channel chan []Contact
	once    sync.Once
}

func CreateChannel() ContactChannel {
	return ContactChannel{channel: make(chan []Contact), once: sync.Once{}}
}

func (ch *ContactChannel) Close() {
	close(ch.channel)
}

func (ch *ContactChannel) Write(contacts []Contact) {
	ch.once.Do(func() {
		ch.channel <- contacts
	})
}

func (ch *ContactChannel) Read() []Contact {
	return <-ch.channel
}
