package d7024e

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"encoding/json"
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
* if connectTo is "none", it will not connect to another node
* Currently only sends a ping to connectTo, because on testing we only want it to puplish itself to one node.
 */
func StartNode(port int, connectTo string) Network {
	me := NewContact(NewRandomKademliaID(), "localhost:"+string(port))
	rt := NewRoutingTable(me)
	network := NewNetwork(3, Kademlia{RT: rt, K: 20})
	go network.Listen("localhost", port)
	message := NewPingMessage(network.kademlia.RT.me.ID)
	if connectTo != "none" {
		SendMessage(connectTo, message)
	}
	return network

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
	ackItem, _ := json.Marshal(item)
	ack := NewFindValueAckMessage(&message.Sender, &message.RPC_ID, &ackItem)
	newAck, _ := MarshallMessage(ack)

	ConnectAndWrite(addr.String(), newAck)
}

/*
STORE message received over network. Sent to kademlia Store.
*/
func (network Network) OnStoreMessageReceived(message *Message, data StoreMessage, addr net.Addr) {
	network.kademlia.Store(data)
	ack := NewStoreAckMessage(&message.Sender, &message.RPC_ID)
	newAck, _ := MarshallMessage(ack)
	//fmt.Println("Sending STORE acknowledge back to ", addr.String(), " with ", newAck)
	ConnectAndWrite(addr.String(), newAck)
}

/*
FIND_VALUE message received over network, sent to kademlia LookupContact.
*/
func (network Network) OnFindNodeMessageReceived(message *Message, data FindNodeMessage, addr net.Addr) {
	target := NewContact(&data.NodeID, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
	contacts := network.kademlia.LookupContact(&target)
	returnMessage := NewFindNodeAckMessage(NewRandomKademliaID(), &message.RPC_ID, &contacts) //TODO: Fix real sender id
	rMsgJson, _ := MarshallMessage(returnMessage)
	//fmt.Println("Sending FIND_NODE acknowledge back to ", addr.String(), " with ", rMsgJson)
	ConnectAndWrite(addr.String(), rMsgJson)
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
	for i := 0; i < network.alpha && i < len(closestContacts); i++ {
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
		} else {
			*counter += 1
			for i := 0; i < network.alpha && i < len(ackMessage.Nodes); i++ {
				go network.FindContactHelper(ackMessage.Nodes[i].Address, message, counter, targetID, ch)
			}
		}
	}
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
	//This is SendFindValueMessage.
}

/*
* Request to find a value over the network.
*/
func (network *Network) SendFindValueMessage(me *KademliaID) Item{
	closest := network.kademlia.RT.FindClosestContacts(me , 3)
	ch := make(chan Item)
	counter := 0

	for i := 0 ; i < network.alpha ; i ++ {
		me := network.kademlia.RT.me
   	message := NewFindValueMessage(me.ID, closest[i].ID)
		go network.FindValueHelper(closest[i].Address, message, &counter, ch)
	}
	item := <-ch
	return item
}

/*
* A helper function for SendFindValueMessage to retreive an item.
*/
func (network *Network) FindValueHelper (addr string, message Message, counter *int, ch chan Item) { //This is correct.
	if *counter >= network.kademlia.K{
		item := Item{}
		ch <- item
		return

	}else{
		_, response, _ := SendMessage(addr, message)
		ack := response.(AckFindValueMessage) //ack Type AckFindValueMessage
		item := Item{}
		err := json.Unmarshal(ack.Value, &item)
		if err != nil {
			return
		}
		if item.Key.Equals(&message.Sender){
			ch <- item
			return
		}else{
			*counter += 1
			for i := 0; i < network.alpha; i++ {
				go network.FindValueHelper(addr, message, counter, ch)
			}
		}
	}
}

/*
Sends a message over the network to the alpha closest neighbors in the routing table and waits for response
from neighbor OnStoreMessageReceived func.
*/
func (network *Network) SendStoreMessage(target *KademliaID, data []byte) []byte {
	//fmt.Println("Testing to send a STORE message")
	closest := network.kademlia.RT.FindClosestContacts(target, network.alpha)
	ch := make(chan []byte)
	counter := 0

	for i := range closest{
		//fmt.Println("Contact [", i ,"], : ", "\n Address : ",closest[i].Address ,"\n ID : ",closest[i].ID, "\n Distance : ", closest[i].distance,"\n")
		me := network.kademlia.RT.me
		message := NewStoreMessage(closest[i].ID, me.ID, &data)
		go network.StoreHelper(closest[i].Address, message, &counter, ch)
	}
	outData := <-ch
	return outData
}

/*
* Helper function for store where a []byte object is received in the response.
*/
func (network *Network) StoreHelper(addr string, message Message, counter *int, ch chan []byte){
	if *counter >= network.alpha{
		data := []byte("")
		ch <- data
		return
	}else{
		rMsg, _, err := SendMessage(addr, message)
		if err != nil{
			return
		}
		if rMsg.Sender.Equals(&message.Sender){
			ch <- []byte("stored")
			return
		} else {
			*counter += 1
			for i := 0; i < network.alpha; i++ {
				go network.StoreHelper(addr, message, counter, ch)
			}
		}
	}
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
