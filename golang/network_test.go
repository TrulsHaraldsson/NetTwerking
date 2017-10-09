package d7024e

import (
	"encoding/json"
	"net"
	"strconv"
	"testing"
	"time"
)

func initKademliaAndNetwork(rt *RoutingTable, port int) (*Kademlia, *Network) {
	network := NewNetwork(3, "localhost:"+strconv.Itoa(port))
	kademlia := Kademlia{rt, 20, &network, nil}
	network.kademlia = &kademlia
	return &kademlia, &network
}

func TestNetworkSendMessage(t *testing.T) {
	go EchoServer(7999)
	time.Sleep(50 * time.Millisecond) //To assure server is up before sending data
	c := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "adress")
	msg := NewPingMessage(&c)
	network := NewNetwork(3, "localhost:3333")
	returnMsg, _, err := network.sendMessage("localhost:7999", msg)
	if err != nil {
		t.Error("SendMessage error ", err)
	}
	if !msg.Equal(returnMsg) {
		t.Error("Message sent is not Equal to Received.", msg, returnMsg)
	}
}

func TestNetworkSendPingMessage(t *testing.T) {
	_, rt := CreateTestRT2()
	_, network := initKademliaAndNetwork(rt, 8019)
	go network.Listen()
	time.Sleep(50 * time.Millisecond)
	_, rt2 := CreateTestRT()
	kademlia2, network2 := initKademliaAndNetwork(rt2, 7)
	pingMsg := NewPingMessage(kademlia2.RT.me)
	msg, err := network2.SendPingMessage("localhost:8019", &pingMsg)
	if err != nil {
		t.Error("error ", err)
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
		//fmt.Println("about to write...")
		if err != nil {
			//fmt.Println("error occured...")
			panic(err)
		} else {
			udpConn.WriteTo(b, addrClient)
		}
		//fmt.Println("written...")
	}
}
