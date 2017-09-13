package d7024e

import (
	"fmt"
	"net"
	"strconv"
	//"strings"
	"encoding/json"
	//"reflect"
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

/*
FIND_NODE message received over network, sent to kademlia LookupData.
*/
func (network Network) OnFindValueMessageReceived(message *Message, data FindValueMessage, addr net.Addr){
	item := network.kademlia.LookupData(&data.ValueID)
	if item.Value != "" {
		ackItem, _ := json.Marshal(item)		
		ack := NewFindValueAckMessage(&message.Sender, &message.RPC_ID, &ackItem)
		newAck, _ := MarshallMessage(ack)
		fmt.Println("Sending FIND_VALUE acknowledge back to sender ", addr.String() ," with item : ", string(newAck))
		ConnectAndWrite(addr.String(), newAck)
	}
}

/*
STORE message received over network. Sent to kademlia Store.
*/
func (network Network) OnStoreMessageReceived(message *Message, data StoreMessage, addr net.Addr) {
		network.kademlia.Store(data)
		ack := NewStoreAckMessage(&message.Sender, &message.RPC_ID)
		newAck, _ := MarshallMessage(ack)
		fmt.Println("Sending STORE acknowledge back to ", addr.String(), " with ", newAck)
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
	fmt.Println("Sending FIND_NODE acknowledge back to ", addr.String(), " with ", rMsgJson)
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
	//This is SendFindValueMessage.
}

func (network *Network) SendFindValueMessage(target *KademliaID) Item{
	fmt.Println("Testing to send a FIND_VALUE message")	
	closest := network.kademlia.RT.FindClosestContacts(target , 3)
	ch := make(chan Item)
	counter := 0
	
	for i := 0 ; i < network.alpha ; i ++ {		
		fmt.Println("Contact [", i ,"], : ", closest[i])
		me := network.kademlia.RT.me
   		message := NewFindValueMessage(me.ID, closest[i].ID)
		go network.FindValueHelper(closest[i].Address, message, &counter, ch)	// This is correct.
		//go network.FindValueHelper(me.ID, closest[i].Address, message, &counter, ch) // For static testing.
	}
	item := <- ch
	return item
}

func (network *Network) FindValueHelper (addr string, message Message, counter *int, ch chan Item) { //This is correct.
//func (network *Network) FindValueHelper (me *KademliaID, addr string, message Message, counter *int, ch chan Item) { // For static testing.	
	fmt.Println("Using FindValueHelper!\n")
	if *counter >= network.kademlia.K{	
		item := Item{}
		fmt.Println("Item should be nothing : ", item.Value ," , ", item.Key) 
		ch <- item
		return
	
	}else{
		//This is correct.
		_, response, _ := SendMessage(addr, message)
		ack := response.(AckFindValueMessage) //ack Type AckFindValueMessage
		item := Item{}
		err := json.Unmarshal(ack.Value, &item)
		if err != nil{
			return
		}
		if item.Key.Equals(&message.Sender){
			ch <- item
			return 
 		//This is correct end.	

		//For static testing.
/*		if me.Equals(&message.Sender){
			item := Item{"Found me!", *me}
			ch <- item
			fmt.Println("Item found, yay!")
			return
*/		//For static testing end.	
		}else{
			*counter += 1
			for i := 0 ; i < network.alpha ; i ++ {
				go network.FindValueHelper(addr, message, counter, ch)
			}
		}
	}
} 

/*
Sends a message over the network to the closest neighbor in the routing table and waits for response 
from neighbor OnStoreMessageReceived func.
*/
func (network *Network) SendStoreMessage(target *KademliaID, data []byte){
	fmt.Println("Testing to send a STORE message")		
	closest := network.kademlia.RT.FindClosestContacts(target , 1)
	for i := range closest{
		fmt.Println("Contact [", i ,"], : ", "\n Address : ",closest[i].Address ,"\n ID : ",closest[i].ID, "\n Distance : ", closest[i].distance,"\n")
	}
	me := network.kademlia.RT.me
	createMessage := NewStoreMessage(closest[0].ID, me.ID, &data)	
	_, _, err := SendMessage(closest[0].Address, createMessage)
	if err != nil{
		fmt.Println("Response is not correct!")
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
