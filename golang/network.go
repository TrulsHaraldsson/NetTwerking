package d7024e

import (
	"fmt"
	"net"
	"strconv"
)

const MESSAGE_SIZE = 1024

type Network struct {
	alpha    int
	kademlia Kademlia
}

func NewNetwork(alpha int, kademlia Kademlia) Network {
	return Network{alpha, kademlia}
}

/*
* Starts a UDP socket listening on port and ip specified.
* When a package is received it will start a new thread handling it.
 */
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
		b, addrClient, err := ReadAnswer(udpConn)
		if err != nil {
			// handle error
			fmt.Println("Error when reading from socket...", err)
		} else {
			m, mData, err2 := UnmarshallMessage(b)
			if err2 != nil {
				fmt.Println("Error when unmarshalling message...", err2)
			} else {
				go network.HandleConnection(m, mData, addrClient)
				fmt.Println("Starting new thread to handle connection...")
			}
		}

	}
}

/* Function to handle an incoming package. Assuming package is a serialized struct of type "Message".
* Depending on the message type, different functions are run.
* Also gives a answer back to the caller.
 */
func (network Network) HandleConnection(message Message, mData interface{}, addr net.Addr) {
	switch message.MsgType {
	case PING:
		fmt.Println("ping")
	case FIND_NODE:
		fmt.Println("Searching for node.")
		network.OnFindNodeMessageReceived(&message, mData.(FindNodeMessage), addr)
	case FIND_VALUE:
		fmt.Println("Searching for value.")
		valueMessage := mData.(FindValueMessage)
		item := network.kademlia.LookupData(&valueMessage.ValueID)
		if item.Value != "" {
			fmt.Println("Item : ", item)
			//ack new find received message back to sender.
			fmt.Println("Sending FIND_VALUE acknowledge back to sender!")
		}else{
			// call closest neighbors if they have value
			fmt.Println("Sending lookup in 3 separate neighbor nodes if they have value")
		}
		//TODO: fix rest
	case STORE:
		fmt.Println("Storing.") //TODO: Put in function like for FIND_NODE above	
		storeMessage := mData.(StoreMessage) //TODO: Send this instead of message.Data below, need to alter kademlia.store to take a correct parameters
		network.kademlia.Store(storeMessage)
		ack := NewStoreAckMessage(&message.Sender, &message.RPC_ID)
		newAck, _ := MarshallMessage(ack)
		ConnectAndWrite(addr.String(), newAck)
		fmt.Println("Sending STORE acknowledge message back!")

	default:
		fmt.Println("Wrong syntax in message, ignoring it...")
	}
}

/*
* Creates a string address
 */
func CreateAddr(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

/*
* When the message received is a FindNodeMessage, this
 */
func (network Network) OnFindNodeMessageReceived(message *Message, data FindNodeMessage, addr net.Addr) {
	fmt.Println("looking up node")
	target := NewContact(&data.NodeID, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
	contacts := network.kademlia.LookupContact(&target)
	returnMessage := NewFindNodeAckMessage(NewRandomKademliaID(), &message.RPC_ID, &contacts) //TODO: Fix real sender id
	rMsgJson, _ := MarshallMessage(returnMessage)
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

/*
* Sends a message of type "Message" to address specified.
* Waits for a response and unmarshalls it as a Message and the MessageData type.
* Returns both Message and MessageData
 */
func SendMessage(addr string, message Message) (Message, interface{}, error) {
	msgJson, err := MarshallMessage(message)
	if err != nil {
		return Message{}, nil, err
	}
	b, err2 := SendData(addr, msgJson)
	if err2 != nil {
		return Message{}, nil, err2
	}
	returnMsg, msgData, err3 := UnmarshallMessage(b)
	return returnMsg, msgData, err3
}

/*
* Sends data to address specified.
* Waits for a response and returns it
* TODO: Don't wait forever...
 */
func SendData(addr string, data []byte) ([]byte, error) {
	var returnMsg []byte
	addrLocal := CreateAddr("localhost", 0)
	addrRemote, _ := net.ResolveUDPAddr("udp", addr)
	udpConn, err := net.ListenPacket("udp", addrLocal)
	fmt.Println("Listening on", udpConn.LocalAddr().String())
	if err != nil {
		return returnMsg, err
	}
	defer udpConn.Close()
	_, err2 := udpConn.WriteTo(data, addrRemote)
	if err2 != nil {
		return returnMsg, err2
	}
	returnMsg, _, err3 := ReadAnswer(udpConn)
	if err3 != nil {
		return returnMsg, err3
	}
	return returnMsg, nil
}

/*
* Reads from the PacketConnection specified,
* returns whats read as a byte array and the adress from where it came.
 */
func ReadAnswer(udpConn net.PacketConn) ([]byte, net.Addr, error) {
	b := make([]byte, MESSAGE_SIZE)
	n, addr, err := udpConn.ReadFrom(b)
	b = b[:n]
	if err != nil {
		return b, addr, err
	}
	return b, addr, nil
}

/*
* Connects to addr and writes the message.
* Does not wait for a response.
 */
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
