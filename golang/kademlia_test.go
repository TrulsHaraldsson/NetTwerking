package d7024e

// According to : (go test -cover -tags KademliaNode) gives 89.6% test coverage atm.

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"
	"fmt"
)

func TestKademliaBootstrap(t *testing.T) {
	k1 := CreateAndStartNode("localhost:11000", "none", nil)
	time.Sleep(time.Millisecond * 50)
	k2 := CreateAndStartNode("localhost:12000", "none", k1.RT.me)
	if k1.RT.Contacts() != 2 || k2.RT.Contacts() != 2 {
		fmt.Println("Contacts : ", k1.RT.Contacts())
		t.Error("Wrong amount of contacts in rt after bootstrap...")
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

	storage := NewStorage()
	kademlia := Kademlia{RT: rt, K: 20, storage: &storage}
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

func TestKademliaSendStoreMessage2(t *testing.T) {
	_, rt := CreateTestRT10()
	_, network := initKademliaAndNetwork(rt, 9500)

	go network.Listen()

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT11()
	_, network2 := initKademliaAndNetwork(rt2, 3)
	filename2 := "filenameX100"
	data2 := []byte("Testing a fucking shit send.")
	network2.kademlia.SendStoreMessage(&filename2, &data2)
	network2.kademlia.DeleteFile(filename2)
}

func TestKademliaSendStoreMessage(t *testing.T) {
	_, rt := CreateTestRT8()
	_, network := initKademliaAndNetwork(rt, 8002)
	go network.Listen()

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT9()
	_, network2 := initKademliaAndNetwork(rt2, 4)

	filename := "filenameX300"
	data := []byte("Testing a fucking shit send.")
	network2.kademlia.SendStoreMessage(&filename, &data)
	fmt.Println("Should be error printed here.")
	network2.kademlia.DeleteFile(filename)

}

/*
* contact searched for is offline, so timeout will occur...
* closest contact found is not the one searched for, since it is offline.
 */
func TestKademliaSendFindContactMessage(t *testing.T) {
	_, rt := CreateTestRT18()
	_, network := initKademliaAndNetwork(rt, 9102)
	go network.Listen()

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT19()
	_, network2 := initKademliaAndNetwork(rt2, 5)

	contact := network2.kademlia.SendFindContactMessage(
		NewKademliaID("1111111100000000000000000000000000000000"))

	if (len(contact) < 1) || (!contact[0].ID.Equals(NewKademliaID("1111111100000000000000000000000000000000"))) {
		t.Error("contacts are not equal", contact[0].ID, NewKademliaID("1111111100000000000000000000000000000000"))
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
func TestKademliaRAMSearch(t *testing.T) {
	filename := "filenameXY"
	data := []byte("This is the content of file filenameXY!")
	kademlia, _ := initKademliaAndNetwork(&RoutingTable{}, 12345)
	kID := NewContact(NewRandomKademliaID(), "adress")
	message := NewStoreMessage(&kID, &filename, &data)
	storeMessage := StoreMessage{}
	json.Unmarshal(message.Data, &storeMessage)
	kademlia.Store(storeMessage)
	file := kademlia.Search(&filename)
	if *file == "" {
		bText := []byte(*file)
		bool := bytes.EqualFold(bText, data)
		if bool == false {
			t.Error("File content do not match!\n")
		}
	}
}

func TestKademliaMemorySearch(t *testing.T) {
	name := "filenameXY"
	filename := []byte(name)
	data := []byte("This is the content of file filenameXY!")
	kademlia, _ := initKademliaAndNetwork(&RoutingTable{}, 12345)
	storage := Storage{}
	storage.Memory(filename, data)
	file := kademlia.Search(&name)
	if *file == "" {
		bText := []byte(*file)
		bool := bytes.EqualFold(bText, data)
		if bool == false {
			t.Error("File content do not match!\n")
		}
	}
	path := "./../newfiles/" + name
	os.Remove(path)
}

func TestKademliaSendFindValue(t *testing.T){
	A := CreateAndStartNode("localhost:5001", "none", nil)
	B := CreateAndStartNode("localhost:5002", "none", A.RT.me)

	filename := "findvaluemessage"
	data := []byte("This is content of findvaluemessage!")
	err := A.SendStoreMessage(&filename, &data)
	if err != nil{
		t.Error("Unsuccessful SendStoreMessage!")
	}
	C := CreateAndStartNode("localhost:5003", "none", B.RT.me)
	file := C.SendFindValueMessage(&filename)
	if file == nil{
		t.Error("File not found!")
	}
	var ffile string
	err3 := json.Unmarshal(file, &ffile)
	if err3 != nil{
		t.Error("unmarshalling failure in find-value test.")
	}
	if (ffile != string(data)) {
		t.Error("Strings of content dont match!")
	}
	path := "./../newfiles/" + filename //<-- check if true.
	os.Remove(path)
}
