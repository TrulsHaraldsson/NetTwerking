package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

const MESSAGE_SIZE = 1024

type Network struct {
	alpha    int
	kademlia Kademlia
}

func (network Network) Listen(ip string, port int) {
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
		b := make([]byte, MESSAGE_SIZE)
		n, addrClient, err := udpConn.ReadFrom(b)
		b = b[:n]
		fmt.Println(addrClient)
		//fmt.Println("received bytes1: ", b)
		if err != nil {
			// handle error
			fmt.Println("Error when reading from socket...\nBetter luck next time.")
			//return nil
		} else {
			go network.HandleConnection(b, addrClient)
			fmt.Println("Starting new thread to handle connection...")
		}

	}
}

func (network Network) HandleConnection(bytes []byte, addr net.Addr) {
	var message Message
	err := json.Unmarshal(bytes, &message)
	fmt.Println("err : ", err)
	contact := NewContact(&message.Sender, addr.String())
	switch message.MsgType {
	case PING:
		fmt.Println("ping")
		network.SendPingMessage(&contact)
	case FIND_NODE:
		fmt.Println("looking up node")
		data := FindNodeMessage{}
		json.Unmarshal(message.Data, data)
		target := NewContact(&data.NodeID, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
		contacts := network.kademlia.LookupContact(&target)
		returnMessage := NewFindNodeAckMessage(NewRandomKademliaID(), &message.RPC_ID, &contacts) //TODO: Fix real sender id
		rMsgJson, _ := json.Marshal(returnMessage)
		ConnectAndWrite(addr.String(), rMsgJson)
	case FIND_VALUE:
		fmt.Println("looking up value")
		//TODO: fix rest
	case STORE:
		fmt.Println("storing data")
		network.kademlia.Store(message.Data)
		//fmt.Println("Message.Data : ", message.Data)
		//fmt.Println("string(Message.Data) : ", string(message.Data))

		var storemessage StoreMessage
		err2 := json.Unmarshal(message.Data, &storemessage)
		if err2 != nil {
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

func (network *Network) SendMessage(addr string, message Message) (Message, error) {
	var returnMsg Message
	addrLocal := CreateAddr("localhost", 0)
	addrRemote, _ := net.ResolveUDPAddr("udp", addr)
	udpConn, err := net.ListenPacket("udp", addrLocal)
	fmt.Println("Listening on", udpConn.LocalAddr().String())
	if err != nil {
		return returnMsg, err
	}
	defer udpConn.Close()
	msgJson, errJson := json.Marshal(message)
	if errJson != nil {
		return returnMsg, errJson
	}
	_, err2 := udpConn.WriteTo(msgJson, addrRemote)
	if err2 != nil {
		return returnMsg, err2
	}
	returnMsg, err3 := network.ReadAnswer(udpConn)
	if err3 != nil {
		return returnMsg, err3
	}
	return returnMsg, nil
}

func (network *Network) ReadAnswer(udpConn net.PacketConn) (Message, error) {
	b := make([]byte, MESSAGE_SIZE)
	n, _, err := udpConn.ReadFrom(b)
	b = b[:n]
	msg := Message{}
	if err != nil {
		return msg, err
	}
	err2 := json.Unmarshal(b, msg)
	if err2 != nil {
		return msg, err2
	}
	return msg, nil
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
