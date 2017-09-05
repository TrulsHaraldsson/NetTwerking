package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

type Network struct {
}

func Listen(ip string, port int) {
	addrServer := CreateAddr(ip, port)
	udpConn, err := net.ListenPacket("udp", addrServer)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Listening to port", port)
	}
	defer udpConn.Close()
	for {
		fmt.Println("Listening to incoming connections...")
		b := make([]byte, 1024)
		n, addrClient, err := udpConn.ReadFrom(b)
		b = b[:n]
		fmt.Println(addrClient)
		//fmt.Println("received bytes1: ", b)
		if err != nil {
			// handle error
			fmt.Println("Error when reading from socket...\nBetter luck next time.")
			//return nil
		} else {
			go HandleConnection(b, addrClient)
			fmt.Println("Starting new thread to handle connection...")
		}

	}
}

func HandleConnection(bytes []byte, addr net.Addr) {
	var message Message
	err := json.Unmarshal(bytes, &message)
	fmt.Println("err : ", err)
	network := Network{}
	contact := NewContact(&message.Sender, addr.String())
	switch message.MsgType {
	case "PING":
		fmt.Println("ping")
		network.SendPingMessage(&contact)
	case "FIND_NODE":
		fmt.Println("looking up node")
	case "FIND_VALUE":
		fmt.Println("looking up value")
		//TODO: fix rest
	case "STORE":
		fmt.Println("storing data")
		kademlia := Kademlia{}
		kademlia.Store(message.Data)
		fmt.Println("List update : ", kademlia.GetList())
		storemessage := StoreMessage{}
		err2 := json.Unmarshal(message.Data, &storemessage)
		if err2 != nil{
			panic(err2)
		}
		
		//fmt.Println("Ack RPC_ID: ", storemessage.RPC_ID)		
		ack := NewStoreAckMessage(&message.Sender, &message.RPC_ID)
		newAck, _ := json.Marshal(ack)
		ConnectAndWrite(addr.String(), newAck)
		fmt.Println("Sent acknowledge message back!")		
		
		
	default:
		fmt.Println("Wrong syntax in message, ignoring it...")
	}
	//fmt.Println("received bytes2: ", bytes)
	fmt.Println("handling connection done.")
}

func CreateAddr(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
		
}

func ConnectAndWrite(addr string, message []byte) error {
	addrLocal := CreateAddr("localhost", 0)
	addrRemote, _ := net.ResolveUDPAddr("udp", addr)
	udpConn, err := net.ListenPacket("udp", addrLocal)
	fmt.Println("Listening on", udpConn.LocalAddr().String())
	if err != nil {
		return err
	}
	defer udpConn.Close()
	_, err2 := udpConn.WriteTo(message, addrRemote)
	if err2 != nil {
		return err2
	} else {
		fmt.Println("message written...")
		return nil
	}
}
