package d7024e


import (
	"encoding/json"
	"fmt"
)


const RPC_ID_LENGTH = 20 //Bytes 
const PING = "PING"
const STORE = "STORE"
const FIND_NODE = "FIND_NODE"
const FIND_VALUE = "FIND_VALUE"

type Message struct {
	msgType string
	sender KademliaID
	data []byte
}

type FindNodeMessage struct {
	RPC_ID KademliaID
	nodeID KademliaID
}

type FindValueMessage struct {
	RPC_ID KademliaID
	valueID KademliaID
}

type PingMessage struct {
	RPC_ID KademliaID
}

type StoreMessage struct {
	RPC_ID KademliaID
	key KademliaID
	data []byte
}

func NewFindValueMessage(sender *KademliaID, valueID *KademliaID) Message {
	var msg = Message{}
	msg.msgType = FIND_VALUE
	msg.sender = *sender

	var findValue = FindValueMessage{*NewRandomKademliaID(), *valueID}
	data, error := json.Marshal(findValue)
	
	if( error != nil ) {
		fmt.Println("Error when creating find value message")
	}
	
	msg.data = data
	return msg
}


func NewFindNodeMessage(sender *KademliaID, nodeID *KademliaID) Message {
	var msg = Message{}
	msg.msgType = FIND_NODE
	msg.sender = *sender
	
	var findNode = FindNodeMessage{*NewRandomKademliaID(), *nodeID}
	data, error := json.Marshal(findNode)
	
	if( error != nil ) {
		fmt.Println("Error when creating find node message")
	}
	
	msg.data = data
	return msg
}

func NewPingMessage(sender *KademliaID) Message {
	var msg = Message{}
	msg.msgType = PING
	msg.sender = *sender
	
	var ping = PingMessage{*NewRandomKademliaID()}
	data, error := json.Marshal(ping)
	
	if( error != nil ) {
		fmt.Println("Error when creating ping message")
	}
	
	msg.data = data
	return msg
}

func NewStoreMessage(sender *KademliaID, key *KademliaID, storeData *[]byte) Message {
	var msg = Message{}
	msg.msgType = STORE
	msg.sender = *sender

	var store = StoreMessage{*NewRandomKademliaID(), *key, *storeData}
	data, error := json.Marshal(store)
	
	if( error != nil ) {
		fmt.Println("Error when creating store message")
	}
	
	msg.data = data
	return msg
}


