package main

import (
	"fmt"
	"testing"
	"time"

	"../golang"
)

/*
Connection pattern : nodes B,C,D are connected with A.
A requests a store on all nodes.
*/
func TestStoreOnce(t *testing.T) {
	A := d7024e.NewKademlia(8400, "2111111400000000000000000000000000000000")
	A.Start()
	time.Sleep(50 * time.Millisecond)

	B := d7024e.NewKademlia(8401, "2111111400000000000000000000000000000001")
	B.Start()
	B.Ping("localhost:8400")
	time.Sleep(50 * time.Millisecond)

	C := d7024e.NewKademlia(8402, "2111111400000000000000000000000000000002")
	C.Start()
	C.Ping("localhost:8400")
	time.Sleep(50 * time.Millisecond)

	D := d7024e.NewKademlia(8403, "2111111400000000000000000000000000000003")
	D.Start()
	D.Ping("localhost:8400")
	time.Sleep(50 * time.Millisecond)

	fmt.Println("All nodes connected")
	contact := A.SendStoreMessage(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), []byte("Test store"))
	if string(contact) != "stored" {
		t.Error("Value not stored!")
	} else {
		fmt.Println("Complete store.")
	}
}

/*
Connection pattern : nodes A and B are connected.
A Sends multiple stores to B.
*/
func TestMultiStore(t *testing.T) {
	A := d7024e.NewKademlia(8410, "2111111400000000000000000000000000000010")
	A.Start()
	time.Sleep(50 * time.Millisecond)

	B := d7024e.NewKademlia(8411, "2111111400000000000000000000000000000011")
	B.Start()
	B.Ping("localhost:8410")
	time.Sleep(50 * time.Millisecond)

	contact1 := A.SendStoreMessage(d7024e.NewKademliaID("2111111400000000000000000000000000000010"), []byte("First"))
	if string(contact1) != "stored" {
		t.Error("Value not stored!", contact1)
	} else {
		fmt.Println("First store complete")
		contact2 := A.SendStoreMessage(d7024e.NewKademliaID("2111111400000000000000000000000000000010"), []byte("Second"))
		if string(contact2) != "stored" {
			t.Error("Value not stored!", contact2)
		} else {
			fmt.Println("Second store complete.")
			contact3 := A.SendStoreMessage(d7024e.NewKademliaID("2111111400000000000000000000000000000010"), []byte("Third"))
			if string(contact3) != "stored" {
				t.Error("Value not stored!", contact3)
			} else {
				fmt.Println("Third store complete.")
			}
		}
	}
}
