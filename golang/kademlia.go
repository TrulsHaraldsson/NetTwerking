package d7024e

import (
	"encoding/json"
	"fmt"
)

var Items []string
//var Items []Item

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
		/*
		ch := make(chan []byte, 3)
		contacts := kademlia.RT.FindClosestContacts(target.ID, kademlia.K)
		
		for i := range ch {
			go func(){
				
			}()
		}
		*/
		fmt.Println("No value found")
	}
}

func (kademlia *Kademlia) Store(data []byte) {
	
	fmt.Println("Entire data : ", string(data))
	
	var m Message
	err := json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println("Error when unmarshalling", err)
	}
	fmt.Println("Type : ", m.MsgType)
	fmt.Println("Sender : ")
	fmt.Println("Data : ")
/*		MsgType string
	Sender  KademliaID
	RPC_ID KademliaID
	Data    []byte
*/
	//item := Item{}
	//newItem := item(m.Key
	Items = append(Items, string(m.Data))
	fmt.Println("Store func complete")
	return
}

func (kademlia *Kademlia) GetList() []string {
	return Items
}
