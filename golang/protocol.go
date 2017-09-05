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
	MsgType string
	Sender  KademliaID
	Data    []byte
}

type FindNodeMessage struct {
	RPC_ID KademliaID
	NodeID KademliaID
}

type FindValueMessage struct {
	RPC_ID  KademliaID
	ValueID KademliaID
}

type PingMessage struct {
	RPC_ID KademliaID
}

type StoreMessage struct {
	RPC_ID KademliaID
	Key    KademliaID
	Data   []byte
}

type AckStoreMessage struct{
	RPC_ID KademliaID
	Data []byte
}

type AckPingMessage struct{
	RPC_ID KademliaID	
}

type AckFindNodeMessage struct{
	RPC_ID KademliaID
	Type string
	Nodes []byte
}

type AckFindValueMessage struct{
	RPC_ID KademliaID
	Type string
	Values []byte
}

func NewFindValueMessage(sender *KademliaID, valueID *KademliaID) Message {
	var msg = Message{}
	msg.MsgType = FIND_VALUE
	msg.Sender = *sender

	var findValue = FindValueMessage{*NewRandomKademliaID(), *valueID}
	data, error := json.Marshal(findValue)

	if error != nil {
		fmt.Println("Error when creating find value message")
	}

	msg.Data = data
	return msg
}

func NewFindNodeMessage(sender *KademliaID, nodeID *KademliaID) Message {
	var msg = Message{}
	msg.MsgType = FIND_NODE
	msg.Sender = *sender

	var findNode = FindNodeMessage{*NewRandomKademliaID(), *nodeID}
	data, error := json.Marshal(findNode)

	if error != nil {
		fmt.Println("Error when creating find node message")
	}

	msg.Data = data
	return msg
}

func NewPingMessage(sender *KademliaID) Message {
	var msg = Message{}
	msg.MsgType = PING
	msg.Sender = *sender

	var ping = PingMessage{*NewRandomKademliaID()}
	data, error := json.Marshal(ping)

	if error != nil {
		fmt.Println("Error when creating ping message")
	}

	msg.Data = data
	return msg
}

func NewStoreMessage(sender *KademliaID, key *KademliaID, storeData *[]byte) Message {
	var msg = Message{}
	msg.MsgType = STORE
	msg.Sender = *sender
	
	var store = StoreMessage{*NewRandomKademliaID(), *key, *storeData}
	data, error := json.Marshal(store)
	if error != nil {
		fmt.Println("Error when creating store message")
	}

	msg.Data = data
	return msg
}

func NewStoreAckMessage(rpc *KademliaID) AckStoreMessage {
	var ack = AckStoreMessage{}
	ack.RPC_ID = *rpc
	ack.Data = []byte("")
	
	return ack
}

func NewPingAckMessage(rpc *KademliaID) AckPingMessage{
	var ack = AckPingMessage{}
	ack.RPC_ID = *rpc
	
	return ack
}

func NewFindNodeAckMessage(rpc *KademliaID, nodes *[]byte) AckFindNodeMessage{
	var ack = AckFindNodeMessage{}
	ack.RPC_ID = *rpc
	ack.Type = FIND_NODE
	ack.Nodes = *nodes
	
	return ack
}

func NewFindValueAckMessage(rpc *KademliaID, values *[]byte) AckFindValueMessage{
	var ack = AckFindValueMessage{}
	ack.RPC_ID = *rpc
	ack.Type = FIND_VALUE
	ack.Values = *values
	
	return ack
}	
