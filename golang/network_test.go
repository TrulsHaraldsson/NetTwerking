package d7024e

//All current tests work according to : go test -run Network

import (
	"encoding/json"
	"net"
	"strconv"
	"testing"
	"time"
	"fmt"
)

func initKademliaAndNetwork(rt *RoutingTable, port int) (*Kademlia, *Network) {
	network := NewNetwork(3, "localhost:"+strconv.Itoa(port))
	kademlia := Kademlia{rt, 20, &network, nil}
	network.kademlia = &kademlia
	return &kademlia, &network
}

func TestNetworkSendMessage(t *testing.T) {
	fmt.Println("Running go test fails because something weird in network_test, however running go test -run Network works... wtf?")

	//fmt.Println("TESTNETWORK SEND PING MESSAGE")
	go EchoServer(7999)
	//fmt.Println("TESTNETWORK 0")

	time.Sleep(50 * time.Millisecond) //To assure server is up before sending data
	c := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "adress")

	//fmt.Println("TESTNETWORK 1")

	msg := NewPingMessage(&c)
	network := NewNetwork(3, "localhost:3333")
	returnMsg, _, err := network.sendMessage("localhost:7999", msg)
	if err != nil {
		t.Error("SendMessage error ", err)
	}

	//fmt.Println("TESTNETWORK 2")

	if !msg.Equal(returnMsg) {
		t.Error("Message sent is not Equal to Received.", msg, returnMsg)
	} else {
		//fmt.Println("Everything went expected, received correct message ", msg, " and returnMsg ", returnMsg)
	}
	//fmt.Println("TESTNETWORK SEND PING MESSAGE END")
}

func TestNetworkSendPingMessage(t *testing.T) {
	//fmt.Println("TESTNETWORK SEND PING MESSAGE AGAIN")
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
	//fmt.Println("TESTNETWORK SEND PING MESSAGE AGAIN END")
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

/*
func TestNetworkListen(t *testing.T) {
	_, rt := CreateTestRT()
	_, network := initKademliaAndNetwork(rt, 8000)

	go network.Listen()
	time.Sleep(50 * time.Millisecond)
	kID := NewContact(NewRandomKademliaID(), "adress")

	//fmt.Println("HELLO2")
	m2 := NewPingMessage(&kID)
	m2Json, _ := json.Marshal(m2)
	err2 := network.connectAndWrite("localhost:8000", m2Json)
	if err2 != nil {
		t.Error("error2:",err2)
	}

	//fmt.Println("HELLO3")
	m3 := NewFindNodeMessage(&kID, NewRandomKademliaID())
	m3Json, _ := json.Marshal(m3)
	err3 := network.connectAndWrite("localhost:8000", m3Json)
	if err3 != nil {
		t.Error("error3:",err3)
	}

	//fmt.Println("HELLO4")
	filename4 := "filenameX444"
	data4 := []byte("hello world!")
	m4 := NewStoreMessage(&kID, &filename4, &data4)
	m4Json, _ := json.Marshal(m4)
	//fmt.Println("HELLO5")

	err4 := network.connectAndWrite("localhost:8000", m4Json)
	if err4 != nil {
		t.Error("error4:",err4)
	}
	//fmt.Println("HELLO6")

	//Create a file in tmp before this test!
	storage := Storage{}
	filename1 := "filenameBAS"
	bytefilename1 := []byte(filename1)
	bytefiletext1 := []byte("Testing content for filenameBAS")
	storage.Memory(bytefilename1, bytefiletext1)

	m1 := NewFindValueMessage(&kID, NewValueID(&filename1))
	m1Json, _ := json.Marshal(m1)
	err1 := network.connectAndWrite("localhost:8000", m1Json)
	if err1 != nil {
		t.Error("error1:", err1)
	}
	//time.Sleep(400 * time.Millisecond) // To assure server receving data before shutdown.
}
*/
