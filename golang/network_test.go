package d7024e

//All current tests work according to : go test -run Network

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

	filename := "filenameX200"
	kID := NewContact(NewRandomKademliaID(), "adress")
	m1 := NewFindValueMessage(&kID, &filename)
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

	filename2 := "filenameX200"
	data := []byte("hello world!")
	m4 := NewStoreMessage(&kID, &filename2, &data)
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

func TestNetworkSendPingMessage(t *testing.T) {
	_, rt := CreateTestRT2()
	_, network := initKademliaAndNetwork(rt)
	go network.Listen("localhost", 8019)

	time.Sleep(50 * time.Millisecond)

	_, rt2 := CreateTestRT()
	kademlia2, network2 := initKademliaAndNetwork(rt2)

	pingMsg := NewPingMessage(kademlia2.RT.me)
	msg, err := network2.SendPingMessage("localhost:8019", &pingMsg)
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
