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
	storemessage := NewStoreMessage(kID, NewRandomKademliaID(),&data)
	kademlia.Store(storemessage.Data) 
}

func TestKademliaNodeLookupData(t *testing.T){
	fmt.Println("Testing lookup data.")
	data := []byte("World peace!")
	kademlia := Kademlia{}
	kID := NewRandomKademliaID()
	storemessage := NewStoreMessage(kID, NewRandomKademliaID(),&data)
	kademlia.Store(storemessage.Data)
	if kademlia.LookupData(kID) == true {
		fmt.Println("Successful lookup!\n")
	}else{
		fmt.Println("Lookup failure!\n")		
	}
}

func TestKademliaNodeLookupDataFail(t *testing.T){
	fmt.Println("Fail testing lookup data.")
	kademlia := Kademlia{}
	kID := NewRandomKademliaID()
	if kademlia.LookupData(kID) == true {
		fmt.Println("Lookup find item, not good should fail!\n")
	}else{
		fmt.Println("Lookup successfully failed!\n")		
	}
}
