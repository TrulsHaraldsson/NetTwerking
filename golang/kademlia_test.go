package d7024e

import (
	"testing"
	"fmt"
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
		fmt.Println(string(i))
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
	data := []byte("hello world!")
	kademlia := Kademlia{}
	kID := NewRandomKademliaID()
	message := NewStoreMessage(kID, NewRandomKademliaID(),&data)
	fmt.Println("Message : ", message)
	kademlia.Store(message.Data)
	fmt.Println("List : ", kademlia.GetList())	
}

func TestKademliaNodeLookupData(t *testing.T){
	data := []byte("hello world!")
	kademlia := Kademlia{}
	kID := NewRandomKademliaID()
	message := NewStoreMessage(kID, NewRandomKademliaID(),&data)
	fmt.Println("Message : ", message)
	kademlia.Store(message.Data)
	fmt.Println("List : ", kademlia.GetList())	
	
	kademlia.LookupData(string(data))
}

func TestKademliaNodeLookupDataFail(t *testing.T){
	data := []byte("hello world!")
	kademlia := Kademlia{}
	kID := NewRandomKademliaID()
	message := NewStoreMessage(kID, NewRandomKademliaID(),&data)
	fmt.Println("Message : ", message)
	kademlia.Store(message.Data)
	fmt.Println("List : ", kademlia.GetList())	
	searchFail := "Hello"
	kademlia.LookupData(searchFail)
}