package d7024e

import (
	"encoding/json"
	"fmt"
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
 * Sends out a maximum of network.kademlia.K RPC's to find the node in the network with id = kademliaID.
 * Returns the contact if it is found(Can be nil, if not found).
 * TODO: Check first if node is found locally
 * TODO: Add functionality for processing if no node closer is found.
 * TODO: Dont check same node multiple times.
 * TODO: Setup a network to test more of its functionality
 */
func (kademlia *Kademlia) SendFindContactMessage(kademliaID *KademliaID) Contact {
	targetID := kademliaID
	target := NewContact(targetID, "DummyAdress")
	closestContacts := kademlia.LookupContact(&target)
	fmt.Println("How many contacts in rt1", len(closestContacts))
	fmt.Println("Closest contact1", closestContacts[0])
	message := NewFindNodeMessage(&kademlia.RT.me, targetID)
	counter := 0
	ch := make(chan Contact)
	for i := 0; i < kademlia.net.alpha && i < len(closestContacts); i++ {
		go kademlia.FindContactHelper(closestContacts[i].Address, message, &counter, targetID, ch)
	}
	contact := <-ch
	return contact
}

func (kademlia *Kademlia) FindContactHelper(addr string, message Message, counter *int,
	targetID *KademliaID, ch chan Contact) {
	if *counter >= kademlia.K {
		ch <- NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "address")
		return
	} else {
		rMessage, ackMessage, err := kademlia.net.SendFindContactMessage(addr, &message) //TODO: dont ignore error
		if err != nil {
			return
		}
		kademlia.RT.AddContact(rMessage.Sender)
		closestContact := ackMessage.Nodes[0]

		fmt.Println("How many contacts in rt", len(ackMessage.Nodes))
		fmt.Println("Closest contact", closestContact)
		if closestContact.ID.Equals(targetID) {
			fmt.Println("Found Contact!!!", closestContact)
			ch <- closestContact
			*counter += kademlia.K
			return
		} else {
			*counter += 1
			for i := 0; i < kademlia.net.alpha && i < len(ackMessage.Nodes); i++ {
				fmt.Println("Sending to", ackMessage.Nodes[i])
				go kademlia.FindContactHelper(ackMessage.Nodes[i].Address,
					message, counter, targetID, ch)
			}
		}
	}
}

/*
 * Request to find a value over the network.
 */
func (kademlia *Kademlia) SendFindValueMessage(me *KademliaID) Item {
	closest := kademlia.RT.FindClosestContacts(me, kademlia.net.alpha)
	ch := make(chan Item)
	counter := 0

	for i := 0; i < kademlia.net.alpha; i++ {
		me := kademlia.RT.me
		message := NewFindValueMessage(&me, me.ID)
		go kademlia.FindValueHelper(closest[i].Address, message, &counter, ch)
	}
	item := <-ch
	return item
}

/*
 * A helper function for SendFindValueMessage to retreive an item.
 */
func (kademlia *Kademlia) FindValueHelper(addr string, message Message, counter *int, ch chan Item) { //This is correct.
	if *counter >= kademlia.K {
		item := Item{}
		ch <- item
		return

	} else {
		_, response, _ := SendMessage(addr, message)
		ack := response.(AckFindValueMessage)
		item := Item{}
		err := json.Unmarshal(ack.Value, &item)
		if err != nil {
			return
		}
		if item.Key.Equals(message.Sender.ID) {
			ch <- item
			*counter += kademlia.K
			return
		} else {
			*counter += 1
			for i := 0; i < kademlia.net.alpha; i++ {
				go kademlia.FindValueHelper(addr, message, counter, ch)
			}
		}
	}
}

/*
 * Sends a message over the network to the alpha closest neighbors in the routing table and waits for response
 * from neighbor OnStoreMessageReceived func.
 */
func (kademlia *Kademlia) SendStoreMessage(me *KademliaID, data []byte) []byte {
	closest := kademlia.RT.FindClosestContacts(me, kademlia.net.alpha)
	ch := make(chan []byte)
	counter := 0

	for i := range closest {
		//fmt.Println("Contact [", i ,"], : ", "\n Address : ",closest[i].Address ,"\n ID : ",closest[i].ID, "\n Distance : ", closest[i].distance,"\n")
		me := kademlia.RT.me
		message := NewStoreMessage(&closest[i], me.ID, &data)
		go kademlia.StoreHelper(closest[i].Address, message, &counter, ch)
	}
	outData := <-ch
	return outData
}

/*
 * Helper function for store where a []byte object is received in the response.
 */
func (kademlia *Kademlia) StoreHelper(addr string, message Message, counter *int, ch chan []byte) {
	if *counter >= kademlia.K {
		data := []byte("")
		ch <- data
		return
	} else {
		rMsg, _, err := SendMessage(addr, message)
		if err != nil {
			return
		}
		if rMsg.Sender.ID.Equals(message.Sender.ID) {
			ch <- []byte("stored")
			*counter += kademlia.K
			return
		} else {
			*counter += 1
			for i := 0; i < kademlia.net.alpha; i++ {
				go kademlia.StoreHelper(addr, message, counter, ch)
			}
		}
	}
}

/*
 * Checks if a certain hash exist in storage, if it does the item is returned of type Item.
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
	ackItem, _ := json.Marshal(item)
	ack := NewFindValueAckMessage(&message.Sender, &message.RPC_ID, &ackItem)
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
