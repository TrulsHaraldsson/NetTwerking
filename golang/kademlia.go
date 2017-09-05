package d7024e

import (
	"encoding/json"
	"fmt"
)

var Items []string

type Kademlia struct {
	RT *RoutingTable
	K  int
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	contacts := kademlia.RT.FindClosestContacts(target.ID, kademlia.K)
	return contacts
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	var m Message
	err := json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println("Error when unmarshalling", err)
	}
	Items = append(Items, string(m.Data))
	fmt.Println("Store func complete")
	fmt.Println("LIST : ", Items)
	return
}

func (kademlia *Kademlia) GetList() []string {
	return Items
}
