package main

import (
	//"fmt"
	"strconv"
	"testing"
	"time"

	"../golang"
)

func TestFindNode1(t *testing.T) {
	// Node B
	B := d7024e.NewKademlia("localhost:8100", "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	B.StartListening()
	time.Sleep(10 * time.Millisecond)

	// Node A
	A := d7024e.NewKademlia("localhost:8101", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	A.StartListening()
	A.Ping("localhost:8100")
	time.Sleep(10 * time.Millisecond)

	// Node C
	C := d7024e.NewKademlia("localhost:8102", "CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")
	C.StartListening()
	C.Ping("localhost:8100")
	time.Sleep(10 * time.Millisecond)

	// NODE D
	D := d7024e.NewKademlia("localhost:8103", "DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	D.StartListening()
	D.Ping("localhost:8102")
	time.Sleep(10 * time.Millisecond)

	//fmt.Println("All nodes connected", A)
	contacts := A.FindContact(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD"))
	//fmt.Println("Closest contact returned:", contacts[0])
	if !contacts[0].ID.Equals(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")) {
		t.Error("Not correct contact returned", contacts[0])
	}
}

func TestFindNode2(t *testing.T) {
	// Node B
	B := d7024e.NewKademlia("localhost:8105", "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	B.StartListening()
	time.Sleep(10 * time.Millisecond)

	count := 30
	port := 8105
	var Node *d7024e.Kademlia
	for i := 0; i < count; i++ {
		Node = d7024e.NewKademlia("localhost:"+strconv.Itoa(port+i+1), "none")
		Node.StartListening()
		connectAddr := d7024e.CreateAddr("localhost", port+i)
		Node.Ping(connectAddr)
		time.Sleep(10 * time.Millisecond)
	}

	contacts := Node.FindContact(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD"))
	//fmt.Println(contacts)
	//fmt.Println("length:", len(contacts))
	if contacts[0].ID.Equals(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")) {
		t.Error("Contact found, when not supposed to")
	}
	if len(contacts) != 20 {
		t.Error("Supposed to be of length 20, but is", len(contacts))
	}
}
