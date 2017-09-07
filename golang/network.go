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

func NewNetwork(alpha int, kademlia Kademlia) Network{
	return Network{alpha, kademlia}
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
		msg, addrClient, err := ReadAnswer(udpConn)
		if err != nil {
			// handle error
			fmt.Println("Error when reading from socket...", err)
		} else {
			go network.HandleConnection(msg, addrClient)
			fmt.Println("Starting new thread to handle connection...")
		}

	}
}

func (network Network) HandleConnection(message Message, addr net.Addr) {
	switch message.MsgType {
	case PING:
		fmt.Println("ping")
	case FIND_NODE:
		network.OnFindNodeMessageReceived(&message, addr)
	case FIND_VALUE:
		valuemessage := FindValueMessage{}
		err2 := json.Unmarshal(message.Data, &valuemessage)
		if err2 != nil{
			fmt.Println("Error : ", err2)
		}
		
		kademlia := Kademlia{}
		if kademlia.LookupData(&valuemessage.ValueID) == false{
			// call closest neighbors if they have value
			fmt.Println("Sending lookup in 3 separate neighbor nodes if they have value")
		}else{
			//ack new find received message back to sender.
			fmt.Println("Sending ack back to sender!")
		}	
		//TODO: fix rest
	case STORE:
		fmt.Println("storing data") //TODO: Put in function like for FIND_NODE above
		kademlia := Kademlia{}
		kademlia.Store(message.Data)
		storemessage := StoreMessage{}
		err2 := json.Unmarshal(message.Data, &storemessage)
		if err2 != nil {
			fmt.Println("Error : ", err2)
		}
		ack := NewStoreAckMessage(&message.Sender, &message.RPC_ID)
		newAck, _ := json.Marshal(ack)
		ConnectAndWrite(addr.String(), newAck)
		fmt.Println("Sent acknowledge message back!")

	default:
		fmt.Println("Wrong syntax in message, ignoring it...")
	}
}

func CreateAddr(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

func (network Network) OnFindNodeMessageReceived(message *Message, addr net.Addr) {
	fmt.Println("looking up node")
	data := FindNodeMessage{}
	json.Unmarshal(message.Data, data)
	target := NewContact(&data.NodeID, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
	contacts := network.kademlia.LookupContact(&target)
	returnMessage := NewFindNodeAckMessage(NewRandomKademliaID(), &message.RPC_ID, &contacts) //TODO: Fix real sender id
	rMsgJson, _ := json.Marshal(returnMessage)
	ConnectAndWrite(addr.String(), rMsgJson)

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

func SendMessage(addr string, message Message) (Message, error) {
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
	returnMsg, _, err3 := ReadAnswer(udpConn)
	if err3 != nil {
		return returnMsg, err3
	}
	return returnMsg, nil
}

func ReadAnswer(udpConn net.PacketConn) (Message, net.Addr, error) {
	b := make([]byte, MESSAGE_SIZE)
	n, addr, err := udpConn.ReadFrom(b)
	b = b[:n]
	msg := Message{}
	if err != nil {
		return msg, addr, err
	}
	err2 := json.Unmarshal(b, &msg)
	if err2 != nil {
		return msg, addr, err2
	}
	return msg, addr, nil
}

func ConnectAndWrite(addr string, message []byte) error {
	addrLocal := CreateAddr("localhost", 0)
	addrRemote, _ := net.ResolveUDPAddr("udp", addr)
	udpConn, err := net.ListenPacket("udp", addrLocal)
	if err != nil {
		return err
	}
	defer udpConn.Close()
	_, err2 := udpConn.WriteTo(message, addrRemote)
	if err2 != nil {
		return err2
	} else {
		return nil
	}
}
