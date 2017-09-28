package d7024e

import (
	"encoding/json"
	"net"
	"regexp"
	"strconv"
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

/*
 * Creates a new Kademlia instance. Initiate a routing table and
 * links the kademlia instance to a network instance.
 * NOTE: This function wont start listening.
 */
func NewKademlia(port int, kID string) *Kademlia {
	var kademliaID *KademliaID
	if kID != "none" {
		kademliaID = NewKademliaID(kID)
	} else {
		kademliaID = NewRandomKademliaID()
	}
	me := NewContact(kademliaID, "localhost:"+strconv.Itoa(port)) //TODO: Should be real ip, not localhost, but works in local tests.
	rt := NewRoutingTable(me)

	// These three rows link kademlia to network and vice versa
	network := NewNetwork(3)
	kademlia := Kademlia{rt, 20, &network}
	network.kademlia = &kademlia

	return &kademlia
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

/*
 * Returns the kademlia.K closest contacts to target.
 */
func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	contacts := kademlia.RT.FindClosestContacts(target.ID, kademlia.K)
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
	if closestContacts[0].ID.Equals(kademliaID) && !closestContacts[0].Equals(kademlia.RT.me) { //If found locally, and not itself.
		return closestContacts
	}
	message := NewFindNodeMessage(&kademlia.RT.me, targetID) // Create message to be sent.

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
		kademlia.RT.AddContact(rMessage.Sender)        // Updating routingtable with new contact seen.
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
 * TODO: Change accordingly to SendFindContactMessage!
 * Request to find a value over the network.
 */

func (kademlia *Kademlia) SendFindValueMessage(kID *KademliaID) []byte {
	target := NewContact(kID, "DummyAdress")
	closestContacts := kademlia.LookupContact(&target)
	if closestContacts[0].ID.Equals(kID) && !closestContacts[0].Equals(kademlia.RT.me) { //If found locally, and not itself.
		return []byte("")
	}
	message := NewFindValueMessage(&kademlia.RT.me, kID) //FindValueMessage
	tempTable := NewContactStateList(kID, kademlia.K) // Creates the temp table
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
		kademlia.RT.AddContact(rMessage.Sender)        // Updating routingtable with new contact seen.
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
func (kademlia *Kademlia) SendStoreMessage(me *KademliaID, data []byte){

	//1: Use SendFindContactMessage to get list of 'k' closest neighbors.
	list := kademlia.SendFindContactMessage(me)
	//2: Filter out the alpha closest out of those 'k' neighbors.
	for _, v:= range list {
		//3: Send out async messages to each of the neighbors without caring about response.
		message := NewStoreMessage(&v, me, &data)
		kademlia.net.SendStoreMessage(v.Address, &message)
	}

	//4: Done.
}

/*
 * Checks if a certain hash exist in storage, if it does the item is returned of type Item.
 * TODO: Remove LookupData and use Storage after patch.
 * TODO: Check same for Store.
 */
func (kademlia *Kademlia) LookupData(hash *KademliaID) Item {
	newItem := Item{}
	for _, v := range Information {
		if v.Key == *hash {
			newItem.Key = v.Key
			newItem.Value = v.Value
		}
	}
	return newItem
}

/*
 * Stores an item of type Item in a list called Information.
 * TODO: Change to the use of Storage.
 * TODO: Update every store/retreive with Storage.
 */
func (kademlia *Kademlia) Store(m StoreMessage) {
	item := Item{string(m.Data), m.Key}
	Information = append(Information, item)
	return
}

/*
 * This method is called by the network module when a PING message is received.
 */
func (kademlia *Kademlia) OnPingMessageReceived(message *Message, addr net.Addr) {
	msgJson := NewPingAckMessage(&kademlia.RT.me, &message.RPC_ID)
	WriteMessage(addr.String(), msgJson)
}

/*
 * This method is called by the network module when a FIND_VALUE message is received.
 */
func (kademlia *Kademlia) OnFindValueMessageReceived(message *Message, data FindValueMessage, addr net.Addr) {
	item := kademlia.LookupData(&data.ValueID)
	var ackItem []byte
	var ackNodes []Contact
	emptyItem := Item{}

	if item.Key == emptyItem.Key{
		target := NewContact(&data.ValueID, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
		ackNodes = kademlia.LookupContact(&target)
	}else{
		ackItem, _ = json.Marshal(item)
	}

	ack := NewFindValueAckMessage(&message.Sender, &message.RPC_ID, &ackItem, &ackNodes)
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
	returnMessage := NewFindNodeAckMessage(&kademlia.RT.me, &message.RPC_ID, &contacts) //TODO: Fix real sender id
	rMsgJson, _ := MarshallMessage(returnMessage)
	//fmt.Println("Sending FIND_NODE acknowledge back to ", addr.String(), " with ", rMsgJson)
	ConnectAndWrite(addr.String(), rMsgJson)
}

/*
 * Will send a Ping message to the given address.
 */
func (kademlia *Kademlia) Ping(addr string) {
	pingMsg := NewPingMessage(&kademlia.RT.me)
	response, error := kademlia.net.SendPingMessage(addr, &pingMsg)
	if error == nil { // No error
		kademlia.RT.AddContact(response.Sender)
	}
}
