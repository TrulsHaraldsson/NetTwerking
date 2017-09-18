package d7024e

/*
see coverage : go test -cover -tags test
*/

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"
)

func initKademliaAndNetwork(rt *RoutingTable) (*Kademlia, *Network) {
	network := NewNetwork(3)
	kademlia := Kademlia{rt, 20, &network}
	network.kademlia = &kademlia
	return &kademlia, &network
}

func TestNetworkListen(t *testing.T) {
	_, rt := CreateTestRT()

	_, network := initKademliaAndNetwork(rt)

	go network.Listen("localhost", 8000)
	time.Sleep(50 * time.Millisecond)

	kID := NewContact(NewRandomKademliaID(), "adress")
	m1 := NewFindValueMessage(&kID, NewRandomKademliaID())
	m1Json, _ := json.Marshal(m1)
	err1 := ConnectAndWrite("localhost:8000", m1Json)
	if err1 != nil {
		t.Error(err1)
	}

	m2 := NewPingMessage(&kID)
	m2Json, _ := json.Marshal(m2)
	err2 := ConnectAndWrite("localhost:8000", m2Json)
	if err2 != nil {
		t.Error(err2)
	}

	m3 := NewFindNodeMessage(&kID, NewRandomKademliaID())
	m3Json, _ := json.Marshal(m3)
	err3 := ConnectAndWrite("localhost:8000", m3Json)
	if err3 != nil {
		t.Error(err3)
	}

	data := []byte("hello world!")
	m4 := NewStoreMessage(&kID, NewRandomKademliaID(), &data)
	m4Json, _ := json.Marshal(m4)

	err4 := ConnectAndWrite("localhost:8000", m4Json)
	if err4 != nil {
		t.Error(err4)
	}

	err5 := ConnectAndWrite("localhost:8000", []byte("Wrong syntax message!"))
	if err5 != nil {
		t.Error(err5)
	}
	time.Sleep(400 * time.Millisecond) // To assure server receving data before shutdown.
}

func TestNetworkSendMessage(t *testing.T) {
	go EchoServer(7999)
	time.Sleep(50 * time.Millisecond) //To assure server is up before sending data
	c := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "adress")
	msg := NewPingMessage(&c)
	returnMsg, _, err := SendMessage("localhost:7999", msg)
	if err != nil {
		t.Error(err)
	}

	if !msg.Equal(returnMsg) {
		t.Error("Message sent is not Equal to Received.", msg, returnMsg)
	} else {
		fmt.Println("Everything went expected, received correct message ", msg, " and returnMsg ", returnMsg)
	}
}

func TestNetworkSendFindValueMessage(t *testing.T) {
	_, rt := CreateTestRT3()
	kademlia, network := initKademliaAndNetwork(rt)
	kID := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	data := []byte("LookThisUp")
	storeMessage := StoreMessage{*kID, data}
	kademlia.Store(storeMessage)

	go network.Listen("localhost", 8007)

	time.Sleep(50 * time.Millisecond)

	item := network.kademlia.SendFindValueMessage(NewKademliaID("FFFFFFFF00000000000000000000000000000000"))

	if string(data) != item.Value {
		t.Error("Couldn't find the stored value.", item)
	} else {
		fmt.Println("Item returned : ", item, "\n String : ", item.Value, "\n Key : ", item.Key)
	}
}

func TestNetworkSendStoreMessage(t *testing.T) {
	_, rt := CreateTestRT4()
	_, network := initKademliaAndNetwork(rt)
	go network.Listen("localhost", 8009)

	time.Sleep(50 * time.Millisecond)

	contact := network.SendStoreMessage(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), []byte("Testing to store"))
	if string(contact) != "stored" {
		t.Error("Store message was not successful.", contact)
	} else {
		fmt.Println("Successful store!", string(contact))
	}
}

func TestNetworkSendFindContactMessage(t *testing.T) {
	contacts, rt := CreateTestRT2()
	_, network := initKademliaAndNetwork(rt)
	go network.Listen("localhost", 8002)

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT()
	_, network2 := initKademliaAndNetwork(rt2)

	contact := network2.kademlia.SendFindContactMessage(NewKademliaID("1111111100000000000000000000000000000000"))
	if !contact.Equals(contacts[0]) {
		t.Error("contacts are not equal", contact, contacts[0])
	}
	contact2 := network2.kademlia.SendFindContactMessage(NewKademliaID("1111111100000000000000000000000000000001"))
	emptyContact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "address")
	if !contact2.Equals(emptyContact) {
		t.Error("Other contact than default found, when not supposed to...", contact2)
	}
}

func TestNetworkSendPingMessage(t *testing.T) {
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

func EchoServer(port int) {
	addrServer := CreateAddr("localhost", port)
	udpConn, err := net.ListenPacket("udp", addrServer)
	if err != nil {
		panic(err)
	}
	defer udpConn.Close()
	for {
		b := make([]byte, MESSAGE_SIZE)
		n, addrClient, err := udpConn.ReadFrom(b)
		b = b[:n]
		msg := Message{}
		json.Unmarshal(b, &msg)
		fmt.Println("about to write...")
		if err != nil {
			fmt.Println("error occured...")
			panic(err)
		} else {
			udpConn.WriteTo(b, addrClient)
		}
		fmt.Println("written...")

	}
}
