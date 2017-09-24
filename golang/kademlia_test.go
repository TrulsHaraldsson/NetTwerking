package d7024e

// According to : (go test -cover -tags KademliaNode) gives 89.6% test coverage atm.

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

func TestKademliaSendFindValueMessage(t *testing.T) {
	_, rt := CreateTestRT3()
	kademlia, network := initKademliaAndNetwork(rt)
	kID := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	data := []byte("LookThisUp")
	storeMessage := StoreMessage{*kID, data}
	kademlia.Store(storeMessage)

	go network.Listen("localhost", 8007)

	time.Sleep(50 * time.Millisecond)

	item := kademlia.SendFindValueMessage(NewKademliaID("FFFFFFFF00000000000000000000000000000000"))

	if string(data) != item.Value {
		t.Error("Couldn't find the stored value.", item)
	} else {
		fmt.Println("Item returned : ", item, "\n String : ", item.Value, "\n Key : ", item.Key)
	}
}

func TestKademliaSendStoreMessage(t *testing.T) {
	_, rt := CreateTestRT4()
	_, network := initKademliaAndNetwork(rt)
	go network.Listen("localhost", 8009)

	time.Sleep(50 * time.Millisecond)

	contact := network.kademlia.SendStoreMessage(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), []byte("Testing to store"))
	if string(contact) != "stored" {
		t.Error("Store message was not successful.", contact)
	} else {
		fmt.Println("Successful store!", string(contact))
	}
}

func TestKademliaSendFindContactMessage(t *testing.T) {
	_, rt := CreateTestRT2()
	_, network := initKademliaAndNetwork(rt)
	go network.Listen("localhost", 8002)

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT9()
	_, network2 := initKademliaAndNetwork(rt2)

	contact := network2.kademlia.SendFindContactMessage(
		NewKademliaID("1111111200000000000000000000000000000000"))
	fmt.Println(contact[0])
	if !contact[0].ID.Equals(NewKademliaID("1111111200000000000000000000000000000000")) {
		t.Error("contacts are not equal", contact[0])
	}
}

func TestKademliaSendPingMessage(t *testing.T) {
	_, rt := CreateTestRT2()
	_, network := initKademliaAndNetwork(rt)
	go network.Listen("localhost", 8003)

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT()
	kademlia2, network2 := initKademliaAndNetwork(rt2)

	pingMsg := NewPingMessage(&kademlia2.RT.me)
	msg, err := network2.SendPingMessage("localhost:8003", &pingMsg)
	if err != nil {
		t.Error(err)
	}
	if msg.MsgType != PING_ACK {
		t.Error("Did not receive an ack for the ping message...")
	}
}

func TestKademliaNodeStore(t *testing.T) {
	fmt.Println("Testing store data.")
	data := []byte("hello world!")
	kademlia := Kademlia{}
	kID := NewContact(NewRandomKademliaID(), "adress")
	message := NewStoreMessage(&kID, NewRandomKademliaID(), &data)
	storeMessage := StoreMessage{}
	json.Unmarshal(message.Data, &storeMessage)
	kademlia.Store(storeMessage)
}

func TestKademliaNodeLookupData(t *testing.T) {
	fmt.Println("Testing to lookup data.")
	data := []byte("hello world!")
	kademlia := Kademlia{}
	kID := NewContact(NewRandomKademliaID(), "adress")
	message := NewStoreMessage(&kID, NewRandomKademliaID(), &data)
	storeMessage := StoreMessage{}
	json.Unmarshal(message.Data, &storeMessage)
	kademlia.Store(storeMessage)
	fmt.Println("Returned Item : ", kademlia.LookupData(kID.ID))
}

func TestKademliaNodeLookupDataFail(t *testing.T) {
	fmt.Println("Fail testing lookup data.")
	kademlia := Kademlia{}
	kID := NewRandomKademliaID()
	fmt.Println("Returned Item : ", kademlia.LookupData(kID))
}
