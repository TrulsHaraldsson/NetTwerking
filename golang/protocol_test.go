package d7024e

import (
	"fmt"
	"testing"
)


func TestNewPingMessage(t *testing.T) {
	var sender = NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	var msg = NewPingMessage(sender)
	fmt.Println(msg.msgType)
	fmt.Println(msg.sender)
	fmt.Println(msg.data)
}

func TestNewStoreMessage(t *testing.T) {
	var sender = NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	var key = NewKademliaID("AAAAAAAA00000000000000000000000000000000")
	var data = []byte("Data to be stored")
	var msg = NewStoreMessage(sender, key, &data)
	fmt.Println(msg.msgType)
	fmt.Println(msg.sender)
	fmt.Println(msg.data)
}

func TestNewFindNodeMessage(t *testing.T) {
	var sender = NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	var nodeID = NewKademliaID("AAAAAAAA00000000000000000000000000000000")
	var msg = NewFindNodeMessage(sender, nodeID)
	fmt.Println(msg.msgType)
	fmt.Println(msg.sender)
	fmt.Println(msg.data)
}

func TestNewFindValueMessage(t *testing.T) {
	var sender = NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	var valueID = NewKademliaID("AAAAAAAA00000000000000000000000000000000")
	var msg = NewFindValueMessage(sender, valueID)
	fmt.Println(msg.msgType)
	fmt.Println(msg.sender)
	fmt.Println(msg.data)
}
