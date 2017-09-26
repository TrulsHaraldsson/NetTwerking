package d7024e

import (
	"fmt"
	"testing"
)

func TestContactAppendUniqueSorted(t *testing.T) {
	list := NewContactStateList(NewRandomKademliaID(), 20)
	contacts1 := make([]Contact, 50)
	for i := 0; i < 50; i++ {
		//fmt.Println(i)
		id := NewRandomKademliaID()
		contact := NewContact(id, "DummyAddress")
		contacts1[i] = contact

	}
	fmt.Println("##########################")
	list.AppendUniqueSorted(contacts1)
	contacts2 := make([]Contact, 50)
	for i2 := 0; i2 < 50; i2++ {
		//fmt.Println(i2)
		id := NewRandomKademliaID()
		contact := NewContact(id, "DummyAddress")
		contacts2[i2] = contact
	}
	list.AppendUniqueSorted(contacts2)
	fmt.Println("len:", len(list.contacts), "Should be 100")
	ok := true
	for i3 := 0; i3 < 99; i3++ {
		if !list.contacts[i3].Less(list.contacts[i3+1]) {
			ok = false
		}
	}
	if !ok {
		t.Error("Not sorted", list.contacts)
	}
}
