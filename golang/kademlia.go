package d7024e

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

/*
* Kademlia is the main object. It is the API seen from outside.
 */
type Kademlia struct {
	RT      *RoutingTable
	K       int
	net     *Network
	storage *Storage
}

/*
 * Creates a new Kademlia instance. Initiate a routing table and
 * links the kademlia instance to a network instance.
 * NOTE: This function wont start listening.
 */
func NewKademlia(addr string, kID string) *Kademlia {
	os.Mkdir("../newfiles/", 0700) //create dir for files.
	rand.Seed(int64(time.Now().Unix()))
	var kademliaID *KademliaID
	if kID != "none" {
		kademliaID = NewKademliaID(kID)
	} else {
		kademliaID = NewRandomKademliaID()
	}
	me := NewContact(kademliaID, addr)
	rt := newRoutingTable(me)
	storage := NewStorage()

	// These three rows link kademlia to network and vice versa
	network := NewNetwork(3, addr)
	kademlia := Kademlia{rt, 20, &network, &storage}
	network.kademlia = &kademlia

	return &kademlia
}

/*
* called when joining a network. A.k.a bootstrap.
* Creates a kademlia instance, starts listening and
* connects to network if initContact is not nil.
 */
func CreateAndStartNode(address string, kID string, initContact *Contact) *Kademlia {
	kademlia := NewKademlia(address, kID)
	kademlia.StartListening()
	if initContact != nil {
		kademlia.RT.update(*initContact)
		kademlia.JoinNetwork()
	}
	return kademlia
}

/*
 * Start listening to the given port
 */
func (kademlia *Kademlia) StartListening() {
	go kademlia.net.Listen()
}

/*
* Tries to connect to network using its current routingtable.
* Max 3 tries, otherwise panic.
* If success, fill up buckets that needs it.
 */
func (kademlia *Kademlia) JoinNetwork() {
	count := 3
	pass := false
	for i := 0; i < count; i++ {
		contacts := kademlia.FindContact(kademlia.RT.me.ID)
		if len(contacts) > 1 { // for success, amount of contacts should be atleast 2 (self and one asked node).
			pass = true
			break
		}
		time.Sleep(2 * time.Second) // sleep 2 seconds and try again.
	}
	if pass {
		for i := 1; i < kademlia.RT.Size()-2; i++ { // -2 since the two buckets at bottom will already be filled to max.
			id := kademlia.RT.getRandomIDForBucket(i)
			go kademlia.FindContact(id)
		}
	} else {
		panic(errors.New("Could not connect to network..."))
	}
}

/*
 * Returns the kademlia.K closest contacts to target.
 */
func (kademlia *Kademlia) LookupContactLocal(target *Contact) []Contact {
	contacts := kademlia.RT.findClosestContacts(target.ID, kademlia.K)
	return contacts
}

/*
 * Sends out FindNode RPC's to find the node in the network with id = kademliaID.
 * Finishes when the k closest nodes are found, and has responded.
 * Returns the K closest contacts found. Closest first in list
 */
func (kademlia *Kademlia) FindContact(kademliaID *KademliaID) []Contact {
	targetID := kademliaID
	target := NewContact(targetID, "DummyAdress")
	closestContacts := kademlia.LookupContactLocal(&target)
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

/*
* Helper function for find contact. the function that is called recursively.
* Sends the actual find contact RPC's
 */
func (kademlia *Kademlia) FindContactHelper(ContactToSendTo Contact, message Message,
	ch *ContactChannel, tempTable *ContactStateList) {
	rMessage, ackMessage, err :=
		kademlia.net.sendFindContactMessage(ContactToSendTo.Address, &message) // Sending RPC, and waiting for response
	if err != nil {
		//If error in message
		fmt.Println("FIND_CONTACT_HELPER:", err)
		tempTable.SetNotQueried(ContactToSendTo) // Set not queried, so others can try again
	} else {
		//If message received correctly
		kademlia.RT.update(rMessage.Sender)            // Updating routingtable with new contact seen.
		tempTable.AppendUniqueSorted(ackMessage.Nodes) // Appends new nodes into tempTable
		tempTable.MarkReceived(ContactToSendTo)        // Mark this contact received.
	}
	if tempTable.Finished() {
		ch.Write(tempTable.GetKClosestContacts()) // Can only be written to once.
	} else { // If not finished, keep sending out RPC's
		for i := 0; i < 1; /*kademlia.net.alpha*/ i++ { // alpha recursive calls to the closest nodes.
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
func (kademlia *Kademlia) FindValue(filename *string) []byte {
	kademliaID := NewValueID(filename) //SHA1 hash
	fileString := kademliaID.String()
	fileContent := kademlia.SearchFileLocal(&fileString)
	if fileContent != nil {
		return []byte(*fileContent)
	}
	fmt.Println("(Not found local, print for demo purpose..)")
	myself := kademlia.RT.me
	cSearch := NewContact(kademliaID, "no address")
	closestContacts := kademlia.LookupContactLocal(&cSearch) //BackHere //TODO: Should search for filename id.

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
	ch2.CloseData()
	return data
}

/*
* Helper function for find value. the function that is called recursively.
* Sends the actual find value RPC's
 */
func (kademlia *Kademlia) FindValueHelper(ContactToSendTo Contact, message Message, ch2 *DataChannel, tempTable *ContactStateList) {
	rMessage, ackMessage, err :=
		kademlia.net.sendFindValueMessage(ContactToSendTo.Address, &message) // Sending RPC, and waiting for response

	if err != nil {
		fmt.Println("FIND_VALUE_HELPER:", err)
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

	if tempTable.Finished() { // If finished
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
* Sending a store message to k closest contacts to the 160 bit hash of filename.
* filename - Filename in plain text e.g. MyFile.txt
 */
func (kademlia *Kademlia) Store(filename *KademliaID, data *[]byte) {
	valueID := filename
	//1: Use FindContact to get list of 'k' closest neighbors.
	contacts := kademlia.FindContact(valueID)
	strValueID := valueID.String()
	message := NewStoreMessage(kademlia.RT.me, &strValueID, data)
	//2: Store it locally before sending out RPC's
	//kademlia.StoreFileLocal(valueID.String(), *data)
	for _, v := range contacts {
		//3: Send out async messages to each of the neighbors without caring about response.
		if v.Equals(*kademlia.RT.me) {
			kademlia.StoreFileLocal(valueID.String(), *data)
		} else {
			go kademlia.net.sendStoreMessage(v.Address, &message)
		}
	}
}

/*
* Searches Ram and Mem for the file specified.
* Returns a *string with the content of file.
* NOTE: Assumes filename is the hash of the real filename.
 */
func (kademlia *Kademlia) SearchFileLocal(filename *string) *string {
	name := []byte(*filename)
	found := kademlia.storage.Search(name)
	if found == nil {
		return nil
	}
	text := string(found.Text)
	strtext := string(text)
	return &strtext
}

/*
* Stores a file locally in ram and mem.
* NOTE: Assumes filename is the hash of the real filename.
 */
func (kademlia *Kademlia) StoreFileLocal(filename string, data []byte) {
	name := []byte(filename)
	ok := kademlia.storage.Store(name, data)
	//kademlia.storage.Store(name, data)
	//TODO: Start purge/republish timer
	ranInt := rand.Intn(60000)
	ranTime := time.Second * (time.Duration(ranInt) / 1000)
	if ok {
		fmt.Println("File received. Purgin and Republishing in:", ranInt/1000+60, "sec.")
		timer := time.AfterFunc(ranTime+time.Second*60, func() { //TODO: dynamic value
			kademlia.storage.deleteTimer(filename)
			kademlia.PurgeAndRepublish(filename)
		})
		kademlia.storage.addTimer(timer, filename)

	} else { // if file did alraedy exist. Update new time for purge/republish
		fmt.Println("Updating timer for:", filename, "New time:", ranInt/1000+60, "sec.")
		kademlia.storage.updateTimer(ranTime+time.Second*60, filename)
	}

}

/*
* Deletes a file from local storage.
 */
func (kademlia *Kademlia) DeleteFileLocal(name string) bool {
	return kademlia.storage.DeleteFile(name)
}

/*
 * This method is called by the network module when a PING message is received.
 * Sends back a ping ack.
 */
func (kademlia *Kademlia) OnPingMessageReceived(message *Message, addr net.Addr) {
	msgJson := NewPingAckMessage(kademlia.RT.me, &message.RPC_ID)
	kademlia.net.WriteMessage(addr.String(), msgJson)
}

/*
 * This method is called by the network module when a FIND_VALUE message is received.
 * Sends back a Findvalue ack.
 */
func (kademlia *Kademlia) OnFindValueMessageReceived(message *Message, fvMessage FindValueMessage, addr net.Addr) {
	filename := fvMessage.Name.String()
	foundFile := kademlia.SearchFileLocal(&filename)
	var ackFile []byte
	var ackNodes []Contact
	if foundFile == nil {
		target := NewContact(&fvMessage.Name, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
		ackNodes = kademlia.LookupContactLocal(&target)
	} else {
		ackFile = []byte(*foundFile)
	}
	ack := NewFindValueAckMessage(&message.Sender, &message.RPC_ID, &ackFile, &ackNodes)
	newAck, _ := MarshallMessage(ack)
	kademlia.net.connectAndWrite(addr.String(), newAck)
}

/*
 * This method is called by the network module when a STORE message is received.
 * Sends back a Store ack.
 */
func (kademlia *Kademlia) OnStoreMessageReceived(message *Message, data StoreMessage, addr net.Addr) {
	kademlia.StoreFileLocal(data.Name, data.Data)
	ack := NewStoreAckMessage(&message.Sender, &message.RPC_ID)
	newAck, _ := MarshallMessage(ack)
	kademlia.net.connectAndWrite(addr.String(), newAck)
}

/*
 * This method is called by the network module when a FIND_NODE message is received.
 * Sends back a find_node ack.
 */
func (kademlia *Kademlia) OnFindNodeMessageReceived(message *Message, data FindNodeMessage, addr net.Addr) {
	target := NewContact(&data.NodeID, "DUMMY ADRESS") // TODO Check if another than dummy adress is needed
	contacts := kademlia.LookupContactLocal(&target)
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
	fmt.Println(error)
	return false
}

/*
* Deletes file locally, and sends store RPC's to the k closest of it.
 */
func (kademlia *Kademlia) PurgeAndRepublish(filename string) {
	content := kademlia.Purge(filename)
	if content != nil {
		kademlia.Republish(filename, []byte(*content))
	}
}

func (kademlia *Kademlia) Purge(filename string) *string {
	fileContent := kademlia.SearchFileLocal(&filename)
	kademlia.DeleteFileLocal(filename)
	return fileContent
}

func (kademlia *Kademlia) Republish(filename string, content []byte) {
	fmt.Println("Purging and Republishing: ", filename)
	valueID := NewKademliaID(filename)
	kademlia.Store(valueID, &content)
}
