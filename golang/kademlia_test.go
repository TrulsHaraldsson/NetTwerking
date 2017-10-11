package d7024e

// According to : (go test -cover -tags KademliaNode) gives 89.6% test coverage atm.

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
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

func TestKademliaNodeLookupContactLocal(t *testing.T) {
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
	contacts := kademlia.LookupContactLocal(&contactsCorrect[0])
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

func TestKademliaStore2(t *testing.T) {
	_, rt := CreateTestRT10()
	_, network := initKademliaAndNetwork(rt, 9500)

	go network.Listen()

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT11()
	_, network2 := initKademliaAndNetwork(rt2, 3)
	filename2 := "filenameX100"
	data2 := []byte("Testing a send2.")
	network2.kademlia.Store(&filename2, &data2)
	time.Sleep(time.Millisecond * 20)
	network2.kademlia.DeleteFileLocal(filename2)
}

func TestKademliaStore(t *testing.T) {
	_, rt := CreateTestRT8()
	_, network := initKademliaAndNetwork(rt, 8002)
	go network.Listen()

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT9()
	k, network2 := initKademliaAndNetwork(rt2, 4)

	filename := "filenameX300"
	fileID := NewValueID(&filename).String()
	data := []byte("Testing a send5.")
	network2.kademlia.Store(&filename, &data)
	file := k.SearchFileLocal(&fileID)
	if string(*file) != string(data) {
		t.Error("wrong content!")
	}
	fmt.Println("Should be error printed here.")
	time.Sleep(time.Millisecond * 20)
	network2.kademlia.DeleteFileLocal(filename)

}

/*
* contact searched for is offline, so timeout will occur...
* closest contact found is not the one searched for, since it is offline.
 */
func TestKademliaFindContact(t *testing.T) {
	_, rt := CreateTestRT18()
	_, network := initKademliaAndNetwork(rt, 9102)
	go network.Listen()

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT19()
	_, network2 := initKademliaAndNetwork(rt2, 5)

	contact := network2.kademlia.FindContact(
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
	kademlia2, _ := initKademliaAndNetwork(rt2, 6)

	ok := kademlia2.Ping("localhost:8003")
	if !ok {
		t.Error("Ping ack was not received correctly")
	}
}

/*
* Store a file by name and content and then search for it.
* Two in One test (Store, SearchFileLocal).
 */
func TestKademliaRAMSearchFileLocal(t *testing.T) {
	filename := "filenameXY"
	data := []byte("This is the content of file filenameXY!")
	kademlia, _ := initKademliaAndNetwork(&RoutingTable{}, 12345)
	kademlia.StoreFileLocal(filename, data)
	file := kademlia.SearchFileLocal(&filename)
	if file != nil {
		bText := []byte(*file)
		bool := bytes.EqualFold(bText, data)
		if bool == false {
			t.Error("File content do not match!\n")
		}
	} else {
		t.Error("No file found...")
	}
}

func TestKademliaMemorySearchFileLocal(t *testing.T) {
	name := "filenameXY"
	filename := []byte(name)
	data := []byte("This is the content of file filenameXY!")
	kademlia, _ := initKademliaAndNetwork(&RoutingTable{}, 12345)
	storage := Storage{}
	storage.Memory(filename, data)
	file := kademlia.SearchFileLocal(&name)
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

func TestKademliaSendFindValue(t *testing.T) {
	A := CreateAndStartNode("localhost:5001", "none", nil)
	B := CreateAndStartNode("localhost:5002", "none", A.RT.me)

	filename := "findvaluemessage"
	data := []byte("This is content of findvaluemessage!")
	A.Store(&filename, &data)
	time.Sleep(50 * time.Millisecond)
	C := CreateAndStartNode("localhost:5003", "none", B.RT.me)
	C.DeleteFileLocal(filename) //Needed to test "not locally found" since items are stored in RAM and Disk, and disk is shared.
	file := C.FindValue(&filename)
	if file == nil {
		t.Error("File not found!")
	}
	ffile := string(file)
	if ffile != string(data) {
		t.Error("Strings of content dont match!")
	}
	time.Sleep(time.Millisecond * 20)
	fileID := NewValueID(&filename)
	path := "../newfiles/" + fileID.String()
	os.Remove(path)
}
