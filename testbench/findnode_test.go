package main

import (
	"fmt"
	"testing"
	"time"

	"../golang"
)

func TestFindNode(t *testing.T) {
	// Node B
	B := d7024e.NewKademlia(8100, "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	B.Start()
	time.Sleep(10 * time.Millisecond)

	// Node A
	A := d7024e.NewKademlia(8101, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	A.Start()
	A.Ping("localhost:8100")
	time.Sleep(10 * time.Millisecond)

	// Node C
	C := d7024e.NewKademlia(8102, "CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")
	C.Start()
	C.Ping("localhost:8100")
	time.Sleep(10 * time.Millisecond)

	// NODE D
	D := d7024e.NewKademlia(8103, "DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	D.Start()
	D.Ping("localhost:8102")
	time.Sleep(10 * time.Millisecond)

	fmt.Println("All nodes connected", A)
	contacts := A.SendFindContactMessage(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD"))
	fmt.Println("Closest contact returned:", contacts[0])
	if !contacts[0].ID.Equals(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")) {
		t.Error("Not correct contact returned", contacts[0])
	}
}

func TestFindNode2(t *testing.T) {
	// Node B
	B := d7024e.NewKademlia(8105, "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	B.Start()
	time.Sleep(10 * time.Millisecond)

	count := 30
	port := 8105
	var Node *d7024e.Kademlia
	for i := 0; i < count; i++ {
		Node = d7024e.NewKademlia(port+i+1, "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
		Node.Start()
		connectAddr := d7024e.CreateAddr("localhost", port+i)
		Node.Ping(connectAddr)
		time.Sleep(10 * time.Millisecond)
	}

	contacts := Node.SendFindContactMessage(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD"))
	fmt.Println(contacts)
	fmt.Println("length:", len(contacts))
	if contacts[0].ID.Equals(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")) {
		t.Error("Contact found, when not supposed to")
	}
	if len(contacts) != 20 {
		t.Error("Supposed to be of length 20, but is", len(contacts))
	}

}
