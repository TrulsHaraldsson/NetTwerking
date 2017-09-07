package d7024e

import (
	"encoding/json"
	"fmt"
)

var Information []Item

type Item struct{
	Value string
	Key KademliaID
}

type Kademlia struct {
	RT *RoutingTable
	K  int
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	contacts := kademlia.RT.FindClosestContacts(target.ID, kademlia.K)
	return contacts
}

func (kademlia *Kademlia) LookupData(hash *KademliaID) bool {
	// TODO
	found := false	
	for _, v := range Information{
		if v.Key == *hash{
			found = true
		}
	}
	if found == true{
		return true
	}else{
		return false
	}
}	

func (kademlia *Kademlia) Store(data []byte){
	var m StoreMessage
	err := json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println("Error when unmarshalling", err)
	}
	item := Item{string(m.Data), m.Key}
	Information = append(Information, item)	
	return 
}
		
func (kademlia *Kademlia) getInformation() []Item {
	return Information
}
