package d7024e

import (
)

var Information []Item

type Item struct {
	Value string
	Key   KademliaID
}

type Kademlia struct {
	RT *RoutingTable
	K  int
}

/*
* Returns the kademlia.K closest contacts to target.
*/
func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	contacts := kademlia.RT.FindClosestContacts(target.ID, kademlia.K)
	return contacts
}

/*
 Checks if a certain hash exist in storage, if it does the item is returned of type Item.
*/
func (kademlia *Kademlia) LookupData(hash *KademliaID) Item {
	newItem := Item{}	
	for _, v := range Information {
		if v.Key == *hash {
			newItem.Key = v.Key
			newItem.Value = v.Value
		}
	}
	return newItem
}

/*
Stores an item of type Item in a list called Information.
*/
func (kademlia *Kademlia) Store(m StoreMessage) {
	item := Item{string(m.Data), m.Key}
	Information = append(Information, item)
	return
}