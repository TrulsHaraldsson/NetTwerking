package d7024e

import (
	"encoding/json"
	"testing"
)

func TestListen(t *testing.T) {
	go Listen("localhost", 8000)
	kID := NewRandomKademliaID()
	m1 := NewFindValueMessage(kID, NewRandomKademliaID())
	m1Json, _ := json.Marshal(m1)
	ConnectAndWrite("localhost:8000", m1Json)

	m2 := NewPingMessage(kID)
	m2Json, _ := json.Marshal(m2)
	ConnectAndWrite("localhost:8000", m2Json)

	m3 := NewFindNodeMessage(kID, NewRandomKademliaID())
	m3Json, _ := json.Marshal(m3)
	ConnectAndWrite("localhost:8000", m3Json)

	data := []byte("hello world!")
	m4 := NewStoreMessage(kID, NewRandomKademliaID(), &data)
	m4Json, _ := json.Marshal(m4)
	ConnectAndWrite("localhost:8000", m4Json)

	ConnectAndWrite("localhost:8000", []byte("Wrong syntax message!"))
}
