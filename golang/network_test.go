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
	err1 := ConnectAndWrite("localhost:8000", m1Json)
	if err1 != nil {
		panic(err1)
	}

	m2 := NewPingMessage(kID)
	m2Json, _ := json.Marshal(m2)
	err2 := ConnectAndWrite("localhost:8000", m2Json)
	if err2 != nil {
		panic(err2)
	}

	m3 := NewFindNodeMessage(kID, NewRandomKademliaID())
	m3Json, _ := json.Marshal(m3)
	err3 := ConnectAndWrite("localhost:8000", m3Json)
	if err3 != nil {
		panic(err3)
	}

	data := []byte("hello world!")
	m4 := NewStoreMessage(kID, NewRandomKademliaID(), &data)
	m4Json, _ := json.Marshal(m4)

	err4 := ConnectAndWrite("localhost:8000", m4Json)
	if err4 != nil {
		panic(err4)
	}

	err5 := ConnectAndWrite("localhost:8000", []byte("Wrong syntax message!"))
	if err5 != nil {
		panic(err5)
	}
}
