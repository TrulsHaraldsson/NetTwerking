package d7024e

// According to : (go test -cover -tags KademliaNode) gives 89.6% test coverage atm.

//All current test work when calling : go test -run Kademlia

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestKademliaNodeLookupContact(t *testing.T) {
	_, rt := CreateTestRT()

	contactsCorrect := []Contact{}
	contactsCorrect = append(contactsCorrect,
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	contactsCorrect = append(contactsCorrect,
		NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))
	contactsCorrect = append(contactsCorrect,
		NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	contactsCorrect = append(contactsCorrect,
		NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	contactsCorrect = append(contactsCorrect,
		NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003"))
	contactsCorrect = append(contactsCorrect,
		NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))

	kademlia := Kademlia{RT: rt, K: 20}
	contacts := kademlia.LookupContact(&contactsCorrect[0])
	fmt.Println(contacts)
	for i, contact := range contacts {
		fmt.Println(" i : ", i, "contact : ", contact)

		if !contact.ID.Equals(contactsCorrect[i].ID) {
			t.Error("Wrong order in contacts")
			fmt.Println(contact.ID, contactsCorrect[i].ID)
		}
		if i > kademlia.K {
			t.Error("Too many contacts returned")
		}
	}
}

func TestKademliaSendFindValueMessage2(t *testing.T){
	_, rt := CreateTestRT10()
	_, network := initKademliaAndNetwork(rt)

	go network.Listen("localhost", 9500)

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT11()
	_, network2 := initKademliaAndNetwork(rt2)

	contact := network2.kademlia.SendFindContactMessage(
		NewKademliaID("1111111100000000000000000000000000000000"))
	fmt.Println("The initial node is : ", rt2.me.ID,"\n")
	for i, j := range contact {
		fmt.Println("[",i,"]", j ,"\n")
	}

	/*
	Send a store message from 1111111100000000000000000000000000000000 to
	FFFFFFFF00000000000000000000000000000000
	*/
	filename2 := "filenameX100"
	data2 := []byte("Testing a fucking shit send.")
	network2.kademlia.SendStoreMessage(&filename2, &data2)
}

func TestKademliaSendStoreMessage(t *testing.T) {
	_, rt := CreateTestRT8()
	_, network := initKademliaAndNetwork(rt)
	go network.Listen("localhost", 8002)

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT9()
	_, network2 := initKademliaAndNetwork(rt2)

	contact := network2.kademlia.SendFindContactMessage(
		NewKademliaID("1111111200000000000000000000000000000000"))
	fmt.Println(contact)
	if !contact[0].ID.Equals(NewKademliaID("1111111100000000000000000000000000000000")) {
		t.Error("contacts are not equal", contact[0].ID, NewKademliaID("1111111200000000000000000000000000000000"))
	}else{
		filename := "filenameX300"
		data := []byte("Testing a fucking shit send.")
		network2.kademlia.SendStoreMessage(&filename, &data)
	}
}

// contact searched for is offline, so timeout will occur...
// closest contact found is not the one searched for, since it is offline.
func TestKademliaSendFindContactMessage(t *testing.T) {
	_, rt := CreateTestRT8()
	_, network := initKademliaAndNetwork(rt)
	go network.Listen("localhost", 9102)

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT9()
	_, network2 := initKademliaAndNetwork(rt2)

	contact := network2.kademlia.SendFindContactMessage(
		NewKademliaID("1111111100000000000000000000000000000000"))

	fmt.Println("Contact list : ",contact,"\n length :", len(contact))
	if (len(contact) < 1) || (!contact[0].ID.Equals(NewKademliaID("1111111100000000000000000000000000000000"))) {
		t.Error("contacts are not equal", contact[0].ID, NewKademliaID("1111111200000000000000000000000000000000"))
	}
}

func TestKademliaSendPingMessage(t *testing.T) {
	_, rt := CreateTestRT2()
	_, network := initKademliaAndNetwork(rt)
	go network.Listen("localhost", 8003)

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT()
	kademlia2, network2 := initKademliaAndNetwork(rt2)

	pingMsg := NewPingMessage(kademlia2.RT.me)
	msg, err := network2.SendPingMessage("localhost:8003", &pingMsg)
	if err != nil {
		t.Error(err)
	}
	if msg.MsgType != PING_ACK {
		t.Error("Did not receive an ack for the ping message...")
	}
}

/*
* Store a file by name and content and then search for it.
* Two in One test (Store, Search).
*/
func TestKademliaNodeSearch(t *testing.T){
	fmt.Println("Testing to search in Storage.Search()!")
	filename := "filenameX200"
	data := []byte("This is the content of file : filenameX200\n")

	kademlia := Kademlia{}
	kID := NewContact(NewRandomKademliaID(), "adress")
	message := NewStoreMessage(&kID, &filename, &data)
	storeMessage := StoreMessage{}
	json.Unmarshal(message.Data, &storeMessage)
	kademlia.Store(storeMessage)
	file := kademlia.Search(&filename)
	fmt.Println("This is the returned file : ", file, "\nSHA-1 hashed name : ", file.Name, "\nContent : ", string(file.Text),"\n")
}
