package d7024e

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const RPC_ID_LENGTH = 20 //Bytes
const PING = "PING"
const STORE = "STORE"
const FIND_NODE = "FIND_NODE"
const FIND_VALUE = "FIND_VALUE"
const PING_ACK = "PING_ACK"
const STORE_ACK = "STORE_ACK"
const FIND_NODE_ACK = "FIND_NODE_ACK"
const FIND_VALUE_ACK = "FIND_VALUE_ACK"

type Message struct {
	MsgType string
	Sender  KademliaID
	RPC_ID  KademliaID
	Data    []byte
}

type FindNodeMessage struct {
	NodeID KademliaID
}

type FindValueMessage struct {
	ValueID KademliaID
}

type PingMessage struct{}

type StoreMessage struct {
	Key  KademliaID
	Data []byte
}

type AckStoreMessage struct{}

type AckPingMessage struct{}

type AckFindNodeMessage struct {
	Nodes []Contact
}

type AckFindValueMessage struct {
	Value []byte
}

func MarshallMessage(msg Message) ([]byte, error) {
	msgJson, err := json.Marshal(msg)
	return msgJson, err
}

/*
* Assumes data is a Message Struct.
* returns the message and messageData unmarshalled too.
 */
func UnmarshallMessage(data []byte) (Message, interface{}, error) {
	m := Message{}

	err1 := json.Unmarshal(data, &m)
	if err1 != nil {
		return m, m, err1
	}
	switch m.MsgType {
	case FIND_NODE:
		mData := FindNodeMessage{}
		err2 := json.Unmarshal(m.Data, &mData)
		return m, mData, err2
	case FIND_NODE_ACK:
		mData := AckFindNodeMessage{}
		err2 := json.Unmarshal(m.Data, &mData)
		return m, mData, err2
	case FIND_VALUE:
		mData := FindValueMessage{}
		err2 := json.Unmarshal(m.Data, &mData)
		return m, mData, err2
	case FIND_VALUE_ACK:
		mData := AckFindValueMessage{}
		err2 := json.Unmarshal(m.Data, &mData)
		return m, mData, err2
	case STORE:
		mData := StoreMessage{}
		err2 := json.Unmarshal(m.Data, &mData)
		return m, mData, err2
	case STORE_ACK:
		mData := AckStoreMessage{}
		err2 := json.Unmarshal(m.Data, &mData)
		return m, mData, err2
	default:
		return m, nil, err1
	}

}

/*
* returns true if the messages are equal.
 */
func (m1 Message) Equal(m2 Message) bool {
	if m1.MsgType != m2.MsgType {
		return false
	} else if m1.Sender != m2.Sender {
		return false
	} else if m1.RPC_ID != m2.RPC_ID {
		return false
	} else if !bytes.Equal(m1.Data, m2.Data) {
		return false
	} else {
		return true
	}
}

func NewFindValueMessage(sender *KademliaID, valueID *KademliaID) Message {
	var msg = Message{}
	msg.MsgType = FIND_VALUE
	msg.Sender = *sender
	msg.RPC_ID = *NewRandomKademliaID()

	var findValue = FindValueMessage{*valueID}
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
	msg.RPC_ID = *NewRandomKademliaID()

	var findNode = FindNodeMessage{*nodeID}
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
	msg.RPC_ID = *NewRandomKademliaID()

	var ping = PingMessage{}
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
	msg.RPC_ID = *NewRandomKademliaID()

	var store = StoreMessage{*key, *storeData}
	data, error := json.Marshal(store)
	if error != nil {
		fmt.Println("Error when creating store message")
	}

	msg.Data = data
	return msg
}

func NewStoreAckMessage(sender *KademliaID, RPC_ID *KademliaID) Message {
	var msg = Message{}
	msg.MsgType = STORE_ACK
	msg.Sender = *sender
	msg.RPC_ID = *RPC_ID

	var ack = AckStoreMessage{}
	data, error := json.Marshal(ack)
	if error != nil {
		fmt.Println("Error when creating store ack message")
	}

	msg.Data = data
	return msg
}

func NewPingAckMessage(sender *KademliaID, RPC_ID *KademliaID) Message {
	var msg = Message{}
	msg.MsgType = PING_ACK
	msg.Sender = *sender
	msg.RPC_ID = *RPC_ID

	var ack = AckPingMessage{}
	data, error := json.Marshal(ack)
	if error != nil {
		fmt.Println("Error when creating ping ack message")
	}

	msg.Data = data
	return msg
}

func NewFindNodeAckMessage(sender *KademliaID, RPC_ID *KademliaID, nodes *[]Contact) Message {
	var msg = Message{}
	msg.MsgType = FIND_NODE_ACK
	msg.Sender = *sender
	msg.RPC_ID = *RPC_ID

	var ack = AckFindNodeMessage{*nodes}
	data, error := json.Marshal(ack)
	if error != nil {
		fmt.Println("Error when creating find node ack message")
	}

	msg.Data = data
	return msg
}

func NewFindValueAckMessage(sender *KademliaID, RPC_ID *KademliaID, value *[]byte) Message {
	var msg = Message{}
	msg.MsgType = FIND_VALUE_ACK
	msg.Sender = *sender
	msg.RPC_ID = *RPC_ID

	var ack = AckFindValueMessage{*value}
	data, error := json.Marshal(ack)
	if error != nil {
		fmt.Println("Error when creating find value ack message")
	}
	msg.Data = data
	return msg
}
