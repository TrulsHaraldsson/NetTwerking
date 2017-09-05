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
	var foundData string
	found := false
	
	for _, v := range Items{
		if v == hash {
			foundData = v
			found = true
		}
	}
	
	if found == true {
		fmt.Println("Found value : ", foundData)
	}else{
		// CHECK OUT NODES 
		fmt.Println("No value found")
	}
}

func (kademlia *Kademlia) Store(data []byte) {
	var m Message
	err := json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println("Error when unmarshalling", err)
	}
	Items = append(Items, string(m.Data))
	fmt.Println("Store func complete")
	return
}

func (kademlia *Kademlia) GetList() []string {
	return Items
}
