package d7024e

import (
	"encoding/json"
	"fmt"
	"reflect"
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

func (kademlia *Kademlia) LookupData(hash *KademliaID){
	// TODO
	found := false
	
	for _, v := range Information{
		if v.Key == *hash{
			fmt.Println("Found item", v.Key)
			found = true
		}
	}
	if found == true{
		fmt.Println("Successful search for item, returning.\n")
		return
	}else{
		fmt.Println("Failed search for item, keep searching.\n")
		//contacts := kademlia.RT.FindClosestContacts(target.ID, kademlia.K)
		return
	}
}	

func receiveAddress () string {
	addr := "www.google.com"
	fmt.Println("Address : ", addr)
	return addr
}

func (kademlia *Kademlia) createChannels(){
	ch := make(chan []byte)
	//Go func
	go func(){
		addr := receiveAddress()
		ch <- []byte(addr)
	}()
	fmt.Println("Channel value : ", ch)
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
/*
func (kademlia *Kademlia) removeInformation(hash KademliaID){
	found := false	
	for i, v := range Information{
		if v.Key == hash {
			Information = append(Information[:i],Information[i + 1:]...)
			found = true
		}
	}
	if found == true {
		fmt.Println("Item Deleted!\n", Information)
	}else{
		fmt.Println("No value found")
	}
}*/