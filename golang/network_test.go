package d7024e

/*
see coverage : go test -cover -tags test
*/

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

//go test -run MultipleChannels etc

func spawn() chan int {
	fmt.Println("Creating channels")
	ch := make(chan int, 2)
	go func() {
		ch <- rand.Intn(100)
	}()
	return ch
}

func TestMultipleChannels(t *testing.T) {
	fmt.Println("Testing Channels")

	ch1 := spawn()
	ch2 := spawn()

	for i := 0; i < 2; i++ {
		select {
		case n := <-ch1:
			fmt.Printf("ch1 : %d\n", n)

		case n := <-ch2:
			fmt.Printf("ch2 : %d\n", n)
		}
	}
}

func TestListen(t *testing.T) {
	_, rt := CreateTestRT()
	network := Network{alpha: 3, kademlia: Kademlia{RT: rt, K: 20}}
	go network.Listen("localhost", 8000)
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
	fmt.Println(string(m3Json))
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
	time.Sleep(200 * time.Millisecond)
}
