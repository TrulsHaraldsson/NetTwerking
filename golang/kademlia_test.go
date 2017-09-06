package d7024e

import (
	"testing"
	"fmt"
	"reflect"
	//"encoding/json"
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
	fmt.Println("kID : ", kID)
	fmt.Println("Type of kID : ", reflect.TypeOf(kID))
	kademlia.LookupData(kID)
}

func TestKademliaNodeLookupDataFail(t *testing.T){
	fmt.Println("Fail testing lookup data.\n")
	data := []byte("World peace!")
	kademlia := Kademlia{}
	kID := NewRandomKademliaID()
	fmt.Println("kID : ", kID)
	fmt.Println("Type of kID : ", reflect.TypeOf(kID))
	kID2 := NewRandomKademliaID()
	storemessage := NewStoreMessage(kID, NewRandomKademliaID(),&data)
	kademlia.Store(storemessage.Data)
	fmt.Println("kID2 : ", kID2)
	fmt.Println("Type of kID2 : ", reflect.TypeOf(kID2))
	kademlia.LookupData(kID2)
}
/*
func TestKademliaNodeRemove(t *testing.T){
	fmt.Println("Testing to remove list item from Information list.")
	data := []byte("Remove Info")
	kademlia := Kademlia{}
	kID := NewRandomKademliaID()
	fmt.Println("kID :: ", kID)
	storemessage := NewStoreMessage(kID, NewRandomKademliaID(),&data)
	kademlia.Store(storemessage.Data)
	
	var m Message
	err := json.Unmarshal(storemessage.Data, &m)
	if err != nil {
		fmt.Println("Error when test unmarshalling", err)
	}
	fmt.Println("sender : ", m.Sender)	
	kademlia.removeInformation(m.Sender)
}*/	


func TestKademliaNodeCreateChannel(t *testing.T){
	fmt.Println("Testing to create channels.")
	kademlia := Kademlia{}
	kademlia.createChannels()
}