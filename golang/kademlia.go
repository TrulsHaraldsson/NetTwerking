package d7024e

import ("fmt"
	"encoding/json")



type Kademlia struct {
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}
var Items []string

func (kademlia *Kademlia) Store(data []byte) {
	var m Message
	err := json.Unmarshal(data, &m)
	if err != nil{
		fmt.Println("Error when unmarshalling", err)
	}	
	Items = append(Items, string(m.Data))
	fmt.Println("Store func complete")

	return 
}

func (kademlia *Kademlia) GetList() []string{
	return Items
}
