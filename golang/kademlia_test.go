package d7024e

// According to : (go test -cover -tags KademliaNode) gives 89.6% test coverage atm.

import (
	"testing"
	"fmt"
	"encoding/json"
)


func TestKademliaNodeLookupContact(t *testing.T) {
	_, rt := CreateTestRT()

	contactsCorrect := []Contact{}
	contactsCorrect = append(contactsCorrect, NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	contactsCorrect = append(contactsCorrect, NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))
	contactsCorrect = append(contactsCorrect, NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	contactsCorrect = append(contactsCorrect, NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	contactsCorrect = append(contactsCorrect, NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003"))
	contactsCorrect = append(contactsCorrect, NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))

	kademlia := Kademlia{RT: rt, K: 20}
	contacts := kademlia.LookupContact(&contactsCorrect[0])
	fmt.Println(contacts)
	for i, contact := range contacts {
		fmt.Println(" i : ", i , "contact : ", contact)
		
		if !contact.ID.Equals(contactsCorrect[i].ID) {
			t.Error("Wrong order in contacts")
			fmt.Println(contact.ID, contactsCorrect[i].ID)
		}
		if i > kademlia.K {
			t.Error("Too many contacts returned")
		}
	}
}

func TestKademliaNodeStore(t *testing.T){
	fmt.Println("Testing store data.")
	data := []byte("hello world!")
	kademlia := Kademlia{}
	kID := NewRandomKademliaID()
	message := NewStoreMessage(kID, NewRandomKademliaID(),&data)
	storeMessage := StoreMessage{}
	json.Unmarshal(message.Data, &storeMessage)
	kademlia.Store(storeMessage) 
}

func TestKademliaNodeLookupData(t *testing.T){
	fmt.Println("Testing store data.")
	data := []byte("hello world!")
	kademlia := Kademlia{}
	kID := NewRandomKademliaID()
	message := NewStoreMessage(kID, NewRandomKademliaID(),&data)
	storeMessage := StoreMessage{}
	json.Unmarshal(message.Data, &storeMessage)
	kademlia.Store(storeMessage) 
	fmt.Println("Returned Item : ", kademlia.LookupData(kID))
}

func TestKademliaNodeLookupDataFail(t *testing.T){
	fmt.Println("Fail testing lookup data.")
	kademlia := Kademlia{}
	kID := NewRandomKademliaID()
	fmt.Println("Returned Item : ", kademlia.LookupData(kID))
}