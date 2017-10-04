package d7024e

import (
	"encoding/json"
	"net"
	"regexp"
	"strconv"
	"fmt"
)

var Information []Item

type Item struct {
	Value string
	Key   KademliaID
}

type Kademlia struct {
	RT  *RoutingTable
	K   int
	net *Network
}

var storage Storage

/*
 * Creates a new Kademlia instance. Initiate a routing table and
 * links the kademlia instance to a network instance.
 * NOTE: This function wont start listening.
 */
func NewKademlia(ip string, port int, kID string) *Kademlia {
	var kademliaID *KademliaID
	if kID != "none" {
		kademliaID = NewKademliaID(kID)
	} else {
		kademliaID = NewRandomKademliaID()
	}
	me := NewContact(kademliaID, ip+":"+strconv.Itoa(port)) //TODO: Should be real ip, not localhost, but works in local tests.
	rt := newRoutingTable(me)

	// These three rows link kademlia to network and vice versa
	network := NewNetwork(3)
	kademlia := Kademlia{rt, 20, &network}
	network.kademlia = &kademlia

	return &kademlia
}

func CreateAndStartNode(ip string, port int, kID string, connecAddr string) *Kademlia {
	kademlia := NewKademlia(ip, port, kID)
	kademlia.Start()
	if connecAddr != "none" {
		kademlia.JoinNetwork(connecAddr)
	}
	return kademlia
}

/*
 * Start listening to the given port
 */
func (kademlia *Kademlia) Start() {
	port, err := strconv.Atoi(regexp.MustCompile(":").Split(kademlia.RT.me.Address, 2)[1]) //Take port and convert to int
	if err != nil {
		panic(err)
	}
	go kademlia.net.Listen("localhost", port)
}

func (kademlia *Kademlia) JoinNetwork(addr string) {
	kademlia.Ping(addr)
	kademlia.SendFindContactMessage(kademlia.RT.me.ID)
	for i := 1; i < kademlia.RT.Size()-2; i++ {
		id := kademlia.RT.getRandomIDForBucket(i)
		go kademlia.SendFindContactMessage(id)
	}
}

/*
 * Returns the kademlia.K closest contacts to target.
 */
func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	contacts := kademlia.RT.findClosestContacts(target.ID, kademlia.K)
	return contacts
}

/*
 * Sends out FindNode RPC's to find the node in the network with id = kademliaID.
 * Finishes when the k closest nodes are found, and has responded.
 * Returns the K closest contacts found. Closest first in list
 * TODO: When no closer node is found, should it send out more RPC's ?
 * TODO: Setup a network to test more of its functionality
 */
func (kademlia *Kademlia) SendFindContactMessage(kademliaID *KademliaID) []Contact {
	targetID := kademliaID
	target := NewContact(targetID, "DummyAdress")
	closestContacts := kademlia.LookupContact(&target)
	if closestContacts[0].ID.Equals(kademliaID) && !closestContacts[0].Equals(*(kademlia.RT.me)) { //If found locally, and not itself.
		return closestContacts
	}
	message := NewFindNodeMessage(kademlia.RT.me, targetID) // Create message to be sent.

	tempTable := NewContactStateList(targetID, kademlia.K) // Creates the temp table
	tempTable.AppendUniqueSorted(closestContacts)
	ch := CreateChannel()                     // Creates a channel that can only be written to once.
	for i := 0; i < kademlia.net.alpha; i++ { // Start with alpha RPC's
		c := tempTable.GetNextToQuery()
		if c != nil { // if nil, there are no current contacts able to query
			go kademlia.FindContactHelper(*c, message, &ch, &tempTable)
		}
	}
	contacts := ch.Read()
	ch.Close()
	return contacts
}

func (kademlia *Kademlia) FindContactHelper(ContactToSendTo Contact, message Message,
	ch *ContactChannel, tempTable *ContactStateList) {
	rMessage, ackMessage, err :=
		kademlia.net.SendFindContactMessage(ContactToSendTo.Address, &message) // Sending RPC, and waiting for response
	if err != nil {
		tempTable.SetNotQueried(ContactToSendTo) // Set not queried, so others can try again
	} else {
		//fmt.Println(ackMessage.Nodes)
		kademlia.RT.update(rMessage.Sender)            // Updating routingtable with new contact seen.
		tempTable.AppendUniqueSorted(ackMessage.Nodes) // Appends new nodes into tempTable
		tempTable.MarkReceived(ContactToSendTo)        // Mark this contact received.
	}
	//fmt.Println(tempTable.contacts)
	if tempTable.Finished() { // If finished,
		ch.Write(tempTable.GetKClosestContacts()) // Can only be written to once.
	} else {
		for i := 0; i < kademlia.net.alpha; i++ { // alpha recursive calls to the closest nodes.
			c := tempTable.GetNextToQuery()
			if c != nil {
				go kademlia.FindContactHelper(*c, message, ch, tempTable)
			}
		}
	}
}

/*
 * Request to find a value over the network.
 */
func (kademlia *Kademlia) SendFindValueMessage(name *string) []byte {
	myself := kademlia.RT.me
	closestContacts := kademlia.LookupContact(myself) //BackHere
	//if closestContacts[0].ID.Equals(myself.ID) && !closestContacts[0].Equals(*kademlia.RT.me) { //If found locally, and not itself.
	if !closestContacts[0].ID.Equals(myself.ID){
		return []byte("")
	}
	message := NewFindValueMessage(myself, name) //FindValueMessage
	tempTable := NewContactStateList(myself.ID, kademlia.K) // Creates the temp table
	tempTable.AppendUniqueSorted(closestContacts)

	ch1 := CreateChannel() //Fix and see if ch2 is required.
	ch2 := CreateDataChannel()                     // Creates a channel that can only be written to once.
	for i := 0; i < kademlia.net.alpha; i++ { // Start with alpha RPC's
		c := tempTable.GetNextToQuery()
		if c != nil { // if nil, there are no current contacts able to query
			go kademlia.FindValueHelper(*c, message, &ch1, &ch2, &tempTable)
		}
	}
	data := ch2.ReadData()
	ch1.Close()
	ch2.CloseData()
	return data
}

func (kademlia *Kademlia) FindValueHelper(ContactToSendTo Contact, message Message, ch1 *ContactChannel, ch2 *DataChannel, tempTable *ContactStateList){
	rMessage, ackMessage, err :=
		kademlia.net.SendFindValueMessage(ContactToSendTo.Address, &message) // Sending RPC, and waiting for response

	if ackMessage.Value != nil {
		ch2.WriteData(ackMessage.Value) // Can only be written to once.
		return
	}

	if err != nil {
		tempTable.SetNotQueried(ContactToSendTo) // Set not queried, so others can try again
	} else {
		//fmt.Println(ackMessage.Nodes)
		kademlia.RT.update(rMessage.Sender)        // Updating routingtable with
		tempTable.AppendUniqueSorted(ackMessage.Nodes) // Appends new nodes into tempTable
		tempTable.MarkReceived(ContactToSendTo)        // Mark this contact received.
	}

	if tempTable.Finished(){ // If finished,
		ch1.Write(tempTable.GetKClosestContacts()) // Can only be written to once.
	} else {
		for i := 0; i < kademlia.net.alpha; i++ { // alpha recursive calls to the closest nodes.
			c := tempTable.GetNextToQuery()
			if c != nil {
				go kademlia.FindValueHelper(*c, message, ch1, ch2, tempTable)
			}
		}
	}
}

/*
* Sending a store message to neighbors.
*/
func (kademlia *Kademlia) SendStoreMessage(filename *string, data *[]byte){
	myself := kademlia.RT.me

	//1: Use SendFindContactMessage to get list of 'k' closest neighbors.
	list := kademlia.SendFindContactMessage(myself.ID)
	//2: Filter out the alpha closest out of those 'k' neighbors.
	for _, v:= range list {
		//3: Send out async messages to each of the neighbors without caring about response.
		message := NewStoreMessage(&v, filename, data)
		kademlia.net.SendStoreMessage(v.Address, &message)
	}
	//4: Done.
}


func (kademlia *Kademlia) Search(filename *string) file {
	name := []byte(*filename)
	found := storage.Search(name)
	return found
}

/*
 * Stores an item of type Item in a list called Information.
 * TODO: Change to the use of Storage.
 * TODO: Update every store/retreive with Storage.
 */
func (kademlia *Kademlia) Store(m StoreMessage) {
	name := []byte(m.Name)
	storage.RAM(name, m.Data)
	return
}

/*
 * This method is called by the network module when a PING message is received.
 */
func (kademlia *Kademlia) OnPingMessageReceived(message *Message, addr net.Addr) {
	msgJson := NewPingAckMessage(kademlia.RT.me, &message.RPC_ID)
	WriteMessage(addr.String(), msgJson)
}

/*
 * This method is called by the network module when a FIND_VALUE message is received.
 */

func (kademlia *Kademlia) OnFindValueMessageReceived(message *Message, fvMessage FindValueMessage, addr net.Addr){
	foundFile := kademlia.Search(&fvMessage.Name)
	var ackFile []byte
	var ackNodes []Contact
	if foundFile.Text == nil{
		target := NewContact(&message.RPC_ID, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
		ackNodes = kademlia.LookupContact(&target)
	}else{
		ackFile, _ = json.Marshal(foundFile)
	}
	ack := NewFindValueAckMessage(&message.Sender, &message.RPC_ID, &ackFile, &ackNodes)
	newAck, _ := MarshallMessage(ack)
	ConnectAndWrite(addr.String(), newAck)
}

/*
 * This method is called by the network module when a STORE message is received.
 */
func (kademlia *Kademlia) OnStoreMessageReceived(message *Message, data StoreMessage, addr net.Addr) {
	kademlia.Store(data)
	ack := NewStoreAckMessage(&message.Sender, &message.RPC_ID)
	newAck, _ := MarshallMessage(ack)
	ConnectAndWrite(addr.String(), newAck)
}

/*
 * This method is called by the network module when a FIND_NODE message is received.
 */
func (kademlia *Kademlia) OnFindNodeMessageReceived(message *Message, data FindNodeMessage, addr net.Addr) {
	target := NewContact(&data.NodeID, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
	contacts := kademlia.LookupContact(&target)
	returnMessage := NewFindNodeAckMessage(kademlia.RT.me, &message.RPC_ID, &contacts) //TODO: Fix real sender id
	rMsgJson, _ := MarshallMessage(returnMessage)
	//fmt.Println("Sending FIND_NODE acknowledge back to ", addr.String(), " with ", rMsgJson)
	ConnectAndWrite(addr.String(), rMsgJson)
}

/*
 * Will send a Ping message to the given address.
 */
func (kademlia *Kademlia) Ping(addr string) {
	pingMsg := NewPingMessage(kademlia.RT.me)
	response, error := kademlia.net.SendPingMessage(addr, &pingMsg)
	fmt.Println("Response:", response)
	if error == nil { // No error
		kademlia.RT.update(response.Sender)
	}
}
