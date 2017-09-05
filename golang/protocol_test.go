package d7024e

import (
	"testing"
)

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
