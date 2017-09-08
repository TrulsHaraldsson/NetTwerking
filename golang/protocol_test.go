package d7024e

import (
	"errors"
	"fmt"
	"testing"
)

func TestProtocolMarshall(t *testing.T) {
	msgData1 := NewKademliaID("fffffffffffffffffffffffffffffffffffffff0")
	msg1 := NewFindNodeMessage(NewKademliaID("ffffffffffffffffffffffffffffffffffffffff"), msgData1)
	err1 := marshallTestHelper(msg1, FindNodeMessage{*msgData1})
	if err1 != nil {
		t.Error(err1)
	}
	msgData2 := []Contact{}
	msgData2 = append(msgData2, NewContact(NewKademliaID("ffffffffffffffffffffffffffffffffffffffff"), "TestAdress"))
	msg2 := NewFindNodeAckMessage(NewKademliaID("ffffffffffffffffffffffffffffffffffffffff"), NewRandomKademliaID(), &msgData2)
	err2 := marshallTestHelper(msg2, AckFindNodeMessage{msgData2})
	if err2 != nil {
		t.Error(err2)
	}
	msgData3 := []byte("Hello World!")
	msg3 := NewFindValueAckMessage(NewKademliaID("ffffffffffffffffffffffffffffffffffffffff"), NewRandomKademliaID(), &msgData3)
	err3 := marshallTestHelper(msg3, AckFindValueMessage{msgData3})
	if err3 != nil {
		t.Error(err3)
	}
	msgData4 := NewKademliaID("fffffffffffffffffffffffffffffffffffffff0")
	msg4 := NewFindValueMessage(NewKademliaID("ffffffffffffffffffffffffffffffffffffffff"), msgData4)
	err4 := marshallTestHelper(msg4, FindNodeMessage{*msgData4})
	if err4 != nil {
		t.Error(err4)
	}
	msgDataKey := NewKademliaID("fffffffffffffffffffffffffffffffffffffff0")
	msgData5 := []byte("Hello World!")
	msg5 := NewStoreMessage(NewKademliaID("ffffffffffffffffffffffffffffffffffffffff"), msgDataKey, &msgData5)
	err5 := marshallTestHelper(msg5, StoreMessage{*msgDataKey, msgData5})
	if err5 != nil {
		t.Error(err5)
	}

}

func marshallTestHelper(msg Message, msgData interface{}) error {
	msgJson, err1 := MarshallMessage(msg)
	if err1 != nil {
		return err1
	}
	msg2, msgData2, err2 := UnmarshallMessage(msgJson)
	if err2 != nil {
		return err2
	}
	if !msg.Equal(msg2) {
		errorMessage := fmt.Sprint("Messages are Not equal", msg, msg2)
		return errors.New(errorMessage)
	}
	if fmt.Sprint(msgData) != fmt.Sprint(msgData2) {
		errorMessage := fmt.Sprint("MessageData are Not equal", msgData, msgData2)
		return errors.New(errorMessage)
	}
	return nil
}

func TestProtocolNewPingMessage(t *testing.T) {
	var sender = NewKademliaID("ffffffffffffffffffffffffffffffffffffffff")
	var msg = NewPingMessage(sender)

	if msg.MsgType != PING {
		t.Error("Expected message type to be", PING, ", got", msg.MsgType)
	}
	if msg.Sender != *sender {
		t.Error("Expected sender to be", *sender, ", got", msg.Sender)
	}
}

func TestProtocolNewStoreMessage(t *testing.T) {
	var sender = NewKademliaID("ffffffffffffffffffffffffffffffffffffffff")
	var key = NewKademliaID("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	var data = []byte("Data to be stored")
	var msg = NewStoreMessage(sender, key, &data)

	if msg.MsgType != STORE {
		t.Error("Expected message type to be", STORE, ", got", msg.MsgType)
	}

	if msg.Sender != *sender {
		t.Error("Expected sender to be", *sender, ", got", msg.Sender)
	}
}

func TestProtocolNewFindNodeMessage(t *testing.T) {
	var sender = NewKademliaID("ffffffff00000000000000000000000000000000")
	var nodeID = NewKademliaID("aaaaaaaa00000000000000000000000000000000")
	var msg = NewFindNodeMessage(sender, nodeID)

	if msg.MsgType != FIND_NODE {
		t.Error("Expected message type to be", FIND_NODE, ", got", msg.MsgType)
	}

	if msg.Sender != *sender {
		t.Error("Expected sender to be", *sender, ", got", msg.Sender)
	}
}

func TestProtocolNewFindValueMessage(t *testing.T) {
	var sender = NewKademliaID("ffffffff00000000000000000000000000000000")
	var valueID = NewKademliaID("aaaaaaaa00000000000000000000000000000000")

	var msg = NewFindValueMessage(sender, valueID)
	if msg.MsgType != FIND_VALUE {
		t.Error("Expected message type to be", FIND_VALUE, ", got", msg.MsgType)
	}

	if msg.Sender != *sender {
		t.Error("Expected sender to be", *sender, ", got", msg.Sender)
	}
}

func TestProtocolNewStoreAckMessage(t *testing.T) {
	var sender = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	var RPC_ID = NewKademliaID("0000000000000000000000000000000000000000")

	var msg = NewStoreAckMessage(sender, RPC_ID)
	if msg.MsgType != STORE_ACK {
		t.Error("Expected message type to be", STORE_ACK, ", got", msg.MsgType)
	}

	if msg.Sender != *sender {
		t.Error("Expected sender to be", *sender, ", got", msg.Sender)
	}

	if msg.RPC_ID != *RPC_ID {
		t.Error("Expected RPC_ID to be", *RPC_ID, ", got", msg.RPC_ID)
	}
}

func TestProtocolNewPingAckMessage(t *testing.T) {
	var sender = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	var RPC_ID = NewKademliaID("0000000000000000000000000000000000000000")

	var msg = NewPingAckMessage(sender, RPC_ID)
	if msg.MsgType != PING_ACK {
		t.Error("Expected message type to be", PING_ACK, ", got", msg.MsgType)
	}

	if msg.Sender != *sender {
		t.Error("Expected sender to be", *sender, ", got", msg.Sender)
	}

	if msg.RPC_ID != *RPC_ID {
		t.Error("Expected RPC_ID to be", *RPC_ID, ", got", msg.RPC_ID)
	}
}

func TestProtocolNewFindNodeAckMessage(t *testing.T) {
	var sender = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	var RPC_ID = NewKademliaID("0000000000000000000000000000000000000000")
	var nodes = []Contact{}

	var msg = NewFindNodeAckMessage(sender, RPC_ID, &nodes)
	if msg.MsgType != FIND_NODE_ACK {
		t.Error("Expected message type to be", FIND_NODE_ACK, ", got", msg.MsgType)
	}

	if msg.Sender != *sender {
		t.Error("Expected sender to be", *sender, ", got", msg.Sender)
	}

	if msg.RPC_ID != *RPC_ID {
		t.Error("Expected RPC_ID to be", *RPC_ID, ", got", msg.RPC_ID)
	}
}

func TestProtocolNewFindValueAckMessage(t *testing.T) {
	var sender = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	var RPC_ID = NewKademliaID("0000000000000000000000000000000000000000")
	var value = []byte("This is data")

	var msg = NewFindValueAckMessage(sender, RPC_ID, &value)
	if msg.MsgType != FIND_VALUE_ACK {
		t.Error("Expected message type to be", FIND_VALUE_ACK, ", got", msg.MsgType)
	}

	if msg.Sender != *sender {
		t.Error("Expected sender to be", *sender, ", got", msg.Sender)
	}

	if msg.RPC_ID != *RPC_ID {
		t.Error("Expected RPC_ID to be", *RPC_ID, ", got", msg.RPC_ID)
	}
}
