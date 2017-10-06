package d7024e

// According to : (go test -cover -tags KademliaNode) gives 89.6% test coverage atm.

//All current test work when calling : go test -run Kademlia

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestKademliaBootstrap(t *testing.T) {
	k1 := CreateAndStartNode("localhost:11000", "none", "none")

	k2 := CreateAndStartNode("localhost:12000", "none", "localhost:11000")
	if k1.RT.Contacts() != 2 {
		t.Error("Expected k1 to have 2 contacts in its routing table, got", k1.RT.Contacts())
	}

	if k2.RT.Contacts() != 2 {
		t.Error("Expected k2 to have 2 contacts in its routing table, got", k2.RT.Contacts())
	} 
}

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
	//fmt.Println(contacts)
	for i, contact := range contacts {
		//fmt.Println(" i : ", i, "contact : ", contact)

		if !contact.ID.Equals(contactsCorrect[i].ID) {
			//fmt.Println(contact.ID, contactsCorrect[i].ID)
			t.Error("Wrong order in contacts")

		}
		if i > kademlia.K {
			t.Error("Too many contacts returned")
		}
	}
}

func TestKademliaSendFindValueMessage(t *testing.T) {
	_, rt := CreateTestRT10()
	_, network := initKademliaAndNetwork(rt, 9500)

	go network.Listen()

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT11()
	_, network2 := initKademliaAndNetwork(rt2, 3)
	filename2 := "filenameX100"
	data2 := []byte("Testing a fucking shit send.")
	network2.kademlia.SendStoreMessage(&filename2, &data2)
}

func TestKademliaSendStoreMessage(t *testing.T) {
	_, rt := CreateTestRT8()
	_, network := initKademliaAndNetwork(rt, 8002)
	go network.Listen()

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT9()
	_, network2 := initKademliaAndNetwork(rt2, 4)

	contact := network2.kademlia.SendFindContactMessage(
		NewKademliaID("1111111200000000000000000000000000000000"))
	if !contact[0].ID.Equals(NewKademliaID("1111111100000000000000000000000000000000")) {
		t.Error("contacts are not equal", contact[0].ID, NewKademliaID("1111111200000000000000000000000000000000"))
	} else {
		filename := "filenameX300"
		data := []byte("Testing a fucking shit send.")
		network2.kademlia.SendStoreMessage(&filename, &data)
	}
}

/*
* contact searched for is offline, so timeout will occur...
* closest contact found is not the one searched for, since it is offline.
 */
func TestKademliaSendFindContactMessage(t *testing.T) {
	_, rt := CreateTestRT8()
	_, network := initKademliaAndNetwork(rt, 9102)
	go network.Listen()

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT9()
	_, network2 := initKademliaAndNetwork(rt2, 5)

	contact := network2.kademlia.SendFindContactMessage(
		NewKademliaID("1111111100000000000000000000000000000000"))

	if (len(contact) < 1) || (!contact[0].ID.Equals(NewKademliaID("1111111100000000000000000000000000000000"))) {
		t.Error("contacts are not equal", contact[0].ID, NewKademliaID("1111111200000000000000000000000000000000"))
	}
}

func TestKademliaSendPingMessage(t *testing.T) {
	_, rt := CreateTestRT2()
	_, network := initKademliaAndNetwork(rt, 8003)
	go network.Listen()

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT()
	kademlia2, network2 := initKademliaAndNetwork(rt2, 6)

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
func TestKademliaNodeSearch(t *testing.T) {
	filename := "filenameXY"
	data := []byte("This is the content of file filenameXY!")
	kademlia := Kademlia{}
	kID := NewContact(NewRandomKademliaID(), "adress")
	message := NewStoreMessage(&kID, &filename, &data)
	storeMessage := StoreMessage{}
	json.Unmarshal(message.Data, &storeMessage)
	kademlia.Store(storeMessage)
	file := kademlia.Search(&filename)
	bool := bytes.EqualFold(file.Text, data)
	if bool == false {
		t.Error("File content do not match!\n", string(data), "\n", string(file.Text), "\n")
	}
}
