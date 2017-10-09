package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

var Information []Item

type Item struct {
	Value string
	Key   KademliaID
}

type Kademlia struct {
	RT      *RoutingTable
	K       int
	net     *Network
	storage *Storage
}

//var storage Storage

/*
 * Creates a new Kademlia instance. Initiate a routing table and
 * links the kademlia instance to a network instance.
 * NOTE: This function wont start listening.
 */
func NewKademlia(addr string, kID string) *Kademlia {

	os.Mkdir("./../newfiles/", 0700)

	var kademliaID *KademliaID
	if kID != "none" {
		kademliaID = NewKademliaID(kID)
	} else {
		kademliaID = NewRandomKademliaID()
	}
	me := NewContact(kademliaID, addr)
	rt := newRoutingTable(me)
	storage := Storage{}

	// These three rows link kademlia to network and vice versa
	network := NewNetwork(3, addr)
	kademlia := Kademlia{rt, 20, &network, &storage}
	network.kademlia = &kademlia

	return &kademlia
}

func CreateAndStartNode(address string, kID string, initContact *Contact) *Kademlia {
	kademlia := NewKademlia(address, kID)
	kademlia.Start()
	if initContact != nil {
		kademlia.RT.update(*initContact)
		kademlia.JoinNetwork()
	}
	return kademlia
}

/*
 * Start listening to the given port
 */
func (kademlia *Kademlia) Start() {
	go kademlia.net.Listen()
}

func (kademlia *Kademlia) JoinNetwork() {
	for {
		// fmt.Println("routingtable size:", kademlia.RT.Contacts())
		// c := NewContact(NewRandomKademliaID(), "address:123")
		// fmt.Println("routingtable contacts:", kademlia.LookupContact(&c))
		contacts := kademlia.SendFindContactMessage(kademlia.RT.me.ID)
		if len(contacts) > 1 {
			fmt.Println("breaking, len is:", len(contacts))
			break
		}
		fmt.Println("not breaking, len is:", len(contacts))
		time.Sleep(2 * time.Second)
	}
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
	message := NewFindNodeMessage(kademlia.RT.me, targetID) // Create message to be sent.

	tempTable := NewContactStateList(targetID, kademlia.K) // Creates the temp table
	tempTable.AppendUniqueSorted(closestContacts)
	tempTable.MarkReceived(*kademlia.RT.me)
	ch := CreateChannel() // Creates a channel that can only be written to once.

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
		kademlia.net.sendFindContactMessage(ContactToSendTo.Address, &message) // Sending RPC, and waiting for response
	if err != nil {
		fmt.Println(err)
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
 * TODO: Fix all todo's
 */
//func (kademlia *Kademlia) SendFindValueMessage(kademliaID *KademliaID) []byte {
func (kademlia *Kademlia) SendFindValueMessage(filename *string) []byte {
	fileContent := kademlia.Search(filename)
	if fileContent != nil {
		fileJson, err := json.Marshal(fileContent)
		if err != nil {
			return fileJson
		}
		return nil
	}
	kademliaID := NewValueID(filename)
	myself := kademlia.RT.me
	cSearch := NewContact(kademliaID, "no address")
	closestContacts := kademlia.LookupContact(&cSearch) //BackHere //TODO: Should search for filename id.

	message := NewFindValueMessage(myself, kademliaID)       //FindValueMessage
	tempTable := NewContactStateList(kademliaID, kademlia.K) // Creates the temp table //TODO: list should be sorted on filename id
	tempTable.AppendUniqueSorted(closestContacts)
	tempTable.MarkReceived(*kademlia.RT.me)

	//Fix and see if ch2 is required. //TODO: c1channel is not necessary i think...
	ch2 := CreateDataChannel()                // Creates a channel that can only be written to once.
	for i := 0; i < kademlia.net.alpha; i++ { // Start with alpha RPC's
		c := tempTable.GetNextToQuery()
		if c != nil { // if nil, there are no current contacts able to query
			go kademlia.FindValueHelper(*c, message, &ch2, &tempTable)
		}
	}
	data := ch2.ReadData()
	//fmt.Println("SendFindValueMessage: After ReadData")

	ch2.CloseData()
	return data
}

func (kademlia *Kademlia) FindValueHelper(ContactToSendTo Contact, message Message, ch2 *DataChannel, tempTable *ContactStateList) {
	rMessage, ackMessage, err :=
		kademlia.net.sendFindValueMessage(ContactToSendTo.Address, &message) // Sending RPC, and waiting for response

	if err != nil { // TODO: Should be checked before ackMessage.value != nil, otherwise crash could occur.
		tempTable.SetNotQueried(ContactToSendTo) // Set not queried, so others can try again
	} else {
		if ackMessage.Value != nil {
			//fmt.Println("FindValueHelper: Found Value!")
			ch2.WriteData(ackMessage.Value) // Can only be written to once.
			return
		}
		//fmt.Println(ackMessage.Nodes)
		kademlia.RT.update(rMessage.Sender)            // Updating routingtable with
		tempTable.AppendUniqueSorted(ackMessage.Nodes) // Appends new nodes into tempTable
		tempTable.MarkReceived(ContactToSendTo)        // Mark this contact received.
	}

	if tempTable.Finished() { // If finished, //TODO: writing to ch1 has no effect, if file is not found, deadlock will occur, should write something like "not found to ch2"
		ch2.WriteData(nil) // Can only be written to once.
	} else {
		for i := 0; i < kademlia.net.alpha; i++ { // alpha recursive calls to the closest nodes.
			c := tempTable.GetNextToQuery()
			if c != nil {
				go kademlia.FindValueHelper(*c, message, ch2, tempTable)
			}
		}
	}
}

/*
* Sending a store message to neighbors.
* filename - Filename in plain text e.g. MyFile.txt
 */
func (kademlia *Kademlia) SendStoreMessage(filename *string, data *[]byte) *KademliaID {
	valueID := NewValueID(filename)

	//1: Use SendFindContactMessage to get list of 'k' closest neighbors.
	contacts := kademlia.SendFindContactMessage(valueID)
	//2: Filter out the alpha closest out of those 'k' neighbors.
	for _, v := range contacts {
		strValueID := valueID.String()
		//3: Send out async messages to each of the neighbors without caring about response.
		message := NewStoreMessage(kademlia.RT.me, &strValueID, data)
		kademlia.net.sendStoreMessage(v.Address, &message)
	}
	return valueID
	//4: Done.
}

func (kademlia *Kademlia) Search(filename *string) *string {
	name := []byte(*filename)
	found := kademlia.storage.Search(name)
	if found == nil {
		return nil
	}
	text := string(found.Text)
	strtext := string(text)
	return &strtext
}

func (kademlia *Kademlia) Store(m StoreMessage) {
	name := []byte(m.Name)
	kademlia.storage.RAM(name, m.Data)
	return
}

/*
 * This method is called by the network module when a PING message is received.
 */
func (kademlia *Kademlia) OnPingMessageReceived(message *Message, addr net.Addr) {
	msgJson := NewPingAckMessage(kademlia.RT.me, &message.RPC_ID)
	kademlia.net.WriteMessage(addr.String(), msgJson)
}

/*
 * This method is called by the network module when a FIND_VALUE message is received.
 */

func (kademlia *Kademlia) OnFindValueMessageReceived(message *Message, fvMessage FindValueMessage, addr net.Addr) {
	filename := fvMessage.Name.String()
	foundFile := kademlia.Search(&filename)
	var ackFile []byte
	var ackNodes []Contact
	if foundFile == nil {
		target := NewContact(&fvMessage.Name, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
		ackNodes = kademlia.LookupContact(&target)
	} else {
		ackFile, _ = json.Marshal(foundFile)
	}
	ack := NewFindValueAckMessage(&message.Sender, &message.RPC_ID, &ackFile, &ackNodes)
	newAck, _ := MarshallMessage(ack)
	kademlia.net.connectAndWrite(addr.String(), newAck)
}

/*
 * This method is called by the network module when a STORE message is received.
 */
func (kademlia *Kademlia) OnStoreMessageReceived(message *Message, data StoreMessage, addr net.Addr) {
	kademlia.Store(data)
	ack := NewStoreAckMessage(&message.Sender, &message.RPC_ID)
	newAck, _ := MarshallMessage(ack)
	kademlia.net.connectAndWrite(addr.String(), newAck)
}

/*
 * This method is called by the network module when a FIND_NODE message is received.
 */
func (kademlia *Kademlia) OnFindNodeMessageReceived(message *Message, data FindNodeMessage, addr net.Addr) {
	target := NewContact(&data.NodeID, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
	contacts := kademlia.LookupContact(&target)
	returnMessage := NewFindNodeAckMessage(kademlia.RT.me, &message.RPC_ID, &contacts) //TODO: Fix real sender id
	rMsgJson, _ := MarshallMessage(returnMessage)
	kademlia.net.connectAndWrite(addr.String(), rMsgJson)
}

/*
 * Will send a Ping message to the given address.
 * Returns True if there was a response, else False.
 */
func (kademlia *Kademlia) Ping(addr string) bool {
	pingMsg := NewPingMessage(kademlia.RT.me)
	response, error := kademlia.net.SendPingMessage(addr, &pingMsg)
	if error == nil { // No error
		kademlia.RT.update(response.Sender)
		return true
	}
	return false
}
