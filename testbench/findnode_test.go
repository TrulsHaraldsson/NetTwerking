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
	contact := A.SendFindContactMessage(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD"))

	if !contact.ID.Equals(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")) {
		t.Error("Not correct contact returned", contact)
	}
}
