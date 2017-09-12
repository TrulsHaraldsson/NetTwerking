package d7024e

import (
	"errors"
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
		fmt.Println("Ping.")
		network.OnPingMessageReceived(&message, addr)
	case FIND_NODE:
		fmt.Println("Searching for node.")
		network.OnFindNodeMessageReceived(&message, mData.(FindNodeMessage), addr)
	case FIND_VALUE:
		fmt.Println("Searching for value.")
		network.OnFindValueMessageReceived(&message, mData.(FindValueMessage), addr)
		//TODO: fix rest
	case STORE:
		fmt.Println("Storing node info.")
		network.OnStoreMessageReceived(&message, mData.(StoreMessage), addr)

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

func (network *Network) OnPingMessageReceived(message *Message, addr net.Addr) {
	msgJson := NewPingAckMessage(network.kademlia.RT.me.ID, &message.RPC_ID)
	WriteMessage(addr.String(), msgJson)
}

/*
FIND_NODE message received over network, sent to kademlia LookupData.
*/
func (network Network) OnFindValueMessageReceived(message *Message, data FindValueMessage, addr net.Addr) {
	item := network.kademlia.LookupData(&data.ValueID)
	if item.Value != "" {
		fmt.Println("Item : ", item)
		fmt.Println("Sending FIND_VALUE acknowledge back to sender!")
	} else {
		fmt.Println("Sending lookup in 3 separate neighbor nodes if they have value")
	}
}

/*
STORE message received over network. Sent to kademlia Store.
*/
func (network Network) OnStoreMessageReceived(message *Message, data StoreMessage, addr net.Addr) {
	network.kademlia.Store(data)
	ack := NewStoreAckMessage(&message.Sender, &message.RPC_ID)
	WriteMessage(addr.String(), ack)
	fmt.Println("Sending STORE acknowledge message back!")
}

/*
FIND_VALUE message received over network, sent to kademlia LookupContact.
*/
func (network Network) OnFindNodeMessageReceived(message *Message, data FindNodeMessage, addr net.Addr) {
	target := NewContact(&data.NodeID, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
	contacts := network.kademlia.LookupContact(&target)
	returnMessage := NewFindNodeAckMessage(network.kademlia.RT.me.ID, &message.RPC_ID, &contacts)
	WriteMessage(addr.String(), returnMessage)
	fmt.Println("Sending back FIND_NODE acknowledge!")
}

/*
* Sends a ping to given address
 */
func (network *Network) SendPingMessage(addr string) (Message, error) {
	msg := NewPingMessage(network.kademlia.RT.me.ID)
	response, _, err := SendMessage(addr, msg)
	if err != nil {
		return Message{}, err
	}
	if msg.RPC_ID == response.RPC_ID {
		return response, nil
	} else {
		return Message{}, errors.New("Wrong RPC_ID returned, it is not from the server, sent to...")
	}

}

/*
* Sends out a maximum of network.kademlia.K RPC's to find the node in the network with id = kademliaID.
* Returns the contact if it is found(Can be nil, if not found).
* TODO: Check first if node is found locally
* TODO: Add functionality for processing if no node closer is found.
* TODO: Dont check same node multiple times.
* TODO: Dont Crash if less than alpha contacts in routingtable
* TODO: Setup a network to test more of its functionality
 */
func (network *Network) SendFindContactMessage(kademliaID *KademliaID) Contact {
	senderID := NewKademliaID("1111111100000000000000000000000000000000")
	targetID := kademliaID
	target := NewContact(targetID, "DummyAdress")
	closestContacts := network.kademlia.LookupContact(&target)
	message := NewFindNodeMessage(senderID, targetID)
	counter := 0
	ch := make(chan Contact)
	for i := 0; i < network.alpha; i++ {
		go network.FindContactHelper(closestContacts[i].Address, message, &counter, targetID, ch)
	}
	contact := <-ch
	return contact
}

func (network *Network) FindContactHelper(addr string, message Message, counter *int, targetID *KademliaID, ch chan Contact) {
	if *counter >= network.kademlia.K {
		ch <- NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "address")
		return
	} else {
		_, response, _ := SendMessage(addr, message) //TODO: dont ignore error
		ackMessage := response.(AckFindNodeMessage)
		closestContact := ackMessage.Nodes[0]
		if closestContact.ID.Equals(targetID) {
			ch <- closestContact
			*counter += network.kademlia.K
			return
		}
		*counter += 1
		for i := 0; i < network.alpha; i++ {
			go network.FindContactHelper(ackMessage.Nodes[i].Address, message, counter, targetID, ch)
		}
	}

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
* Writes a message of type Message to addr, does not wait for response.
 */
func WriteMessage(addr string, message Message) error {
	msgJson, err := MarshallMessage(message)
	if err != nil {
		return err
	}
	err2 := ConnectAndWrite(addr, msgJson)
	return err2
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
