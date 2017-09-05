package d7024e

import (
	"fmt"
	"testing"
)

func TestLookupContact(t *testing.T) {
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

func spawnCh() chan string{	
	fmt.Println("Creating channels")
	ch := make(chan string,2)
	go func(){
		ch <- "hello"
	}()
	return ch
}

func TestStoreItems(t *testing.T){
	fmt.Println("Testing multiple stores")

	ch1 := spawnCh()
	ch2 := spawnCh()
	
	for i := 0; i < 2; i++{
		select {
			case n := <- ch1: 
				fmt.Printf("ch1 : %s\n", n)
				data := []byte(n)
				
				var kademlia Kademlia
				kademlia.Store(data)
				kademlia.GetList()		
				
				
			case n := <- ch2:
				fmt.Printf("ch2 : %s\n", n)
				data := []byte(n)
				
				var kademlia Kademlia
				kademlia.Store(data)
				kademlia.GetList()
		}
	}
}
