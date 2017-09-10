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

func TestNetworkListen(t *testing.T) {
	_, rt := CreateTestRT()
	network := Network{alpha: 3, kademlia: Kademlia{RT: rt, K: 20}}
	go network.Listen("localhost", 8000)

	kID := NewRandomKademliaID()

	m1 := NewFindValueMessage(kID, NewRandomKademliaID())
	m1Json, _ := json.Marshal(m1)
	err1 := ConnectAndWrite("localhost:8000", m1Json)
	if err1 != nil {
		t.Error(err1)
	}

	m2 := NewPingMessage(kID)
	m2Json, _ := json.Marshal(m2)
	err2 := ConnectAndWrite("localhost:8000", m2Json)
	if err2 != nil {
		t.Error(err2)
	}

	m3 := NewFindNodeMessage(kID, NewRandomKademliaID())
	m3Json, _ := json.Marshal(m3)
	err3 := ConnectAndWrite("localhost:8000", m3Json)
	if err3 != nil {
		t.Error(err3)
	}

	data := []byte("hello world!")
	m4 := NewStoreMessage(kID, NewRandomKademliaID(), &data)
	m4Json, _ := json.Marshal(m4)

	err4 := ConnectAndWrite("localhost:8000", m4Json)
	if err4 != nil {
		t.Error(err4)
	}
	
/*
	err5 := ConnectAndWrite("localhost:8000", []byte("Wrong syntax message!"))
	if err5 != nil {
		t.Error(err5)
	}
*/	time.Sleep(400 * time.Millisecond) // To assure server receving data before shutdown.
}

func TestNetworkSendMessage(t *testing.T) {
	go EchoServer(7999)
	time.Sleep(100 * time.Millisecond) //To assure server is up before sending data
	msg := NewPingMessage(NewKademliaID("FFFFFFFF00000000000000000000000000000000"))
	returnMsg, _, err := SendMessage("localhost:7999", msg)
	if err != nil {
		t.Error(err)
	}

	if !msg.Equal(returnMsg) {
		t.Error("Message sent is not Equal to Received.", msg, returnMsg)
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
