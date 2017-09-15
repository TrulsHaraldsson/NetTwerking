package d7024e

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

const MESSAGE_SIZE = 1024

type Network struct {
	alpha    int
	kademlia *Kademlia
}

func NewNetwork(alpha int) Network {
	return Network{alpha, nil}
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
	case STORE:
		fmt.Println("Storing node info.")
		network.OnStoreMessageReceived(&message, mData.(StoreMessage), addr)

	default:
		fmt.Println("Wrong syntax in message, ignoring it...")
		return
	}
	network.kademlia.RT.AddContact(message.Sender)
}

/*
* Creates a string address
 */
func CreateAddr(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

func (network *Network) OnPingMessageReceived(message *Message, addr net.Addr) {
	network.kademlia.OnPingMessageReceived(message, addr)
}

/*
FIND_NODE message received over network, sent to kademlia LookupData.
*/
func (network Network) OnFindValueMessageReceived(message *Message, data FindValueMessage, addr net.Addr) {
	network.kademlia.OnFindValueMessageReceived(message, data, addr)
}

/*
STORE message received over network. Sent to kademlia Store.
*/
func (network Network) OnStoreMessageReceived(message *Message, data StoreMessage, addr net.Addr) {
	network.kademlia.OnStoreMessageReceived(message, data, addr)
}

/*
FIND_VALUE message received over network, sent to kademlia LookupContact.
*/
func (network Network) OnFindNodeMessageReceived(message *Message, data FindNodeMessage, addr net.Addr) {
	network.kademlia.OnFindNodeMessageReceived(message, data, addr)
}

/*
 * Sends a ping to given address
 */
func (network *Network) SendPingMessage(addr string, msg *Message) (Message, error) {
	//msg := NewPingMessage(&network.kademlia.RT.me)
	response, _, err := SendMessage(addr, *msg)
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
* TODO: Setup a network to test more of its functionality
 */
func (network *Network) SendFindContactMessage(kademliaID *KademliaID) Contact {
	return network.kademlia.SendFindContactMessage(kademliaID)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
	//This is SendFindValueMessage.
}

/*
* Request to find a value over the network.
 */
func (network *Network) SendFindValueMessage(me *KademliaID) Item {
	return network.kademlia.SendFindValueMessage(me)
}

/*
Sends a message over the network to the alpha closest neighbors in the routing table and waits for response
from neighbor OnStoreMessageReceived func.
*/
func (network *Network) SendStoreMessage(me *KademliaID, data []byte) []byte {
	return network.kademlia.SendStoreMEssage(me, data)
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
 */
func SendData(addr string, data []byte) ([]byte, error) {
	var returnMsg []byte
	addrLocal := CreateAddr("localhost", 0)
	addrRemote, _ := net.ResolveUDPAddr("udp", addr)
	udpConn, err := net.ListenPacket("udp", addrLocal)
	udpConn.SetDeadline(time.Now().Add(3 * time.Second)) // Waits for 3 seconds
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
