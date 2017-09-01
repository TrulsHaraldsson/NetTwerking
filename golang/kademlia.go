package d7024e

import ("fmt"
	"encoding/json"
)

type Kademlia struct {
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO: Remember that the value is stored but not globally or anything so that it can't be retreived for later use by other func.
	mess := StoreMessage{}
	err := json.Unmarshal(data,&mess)
	if err != nil{
		fmt.Println("message: ",err)
	}
}
