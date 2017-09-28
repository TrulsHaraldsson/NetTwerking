package main

import (
	"fmt"
	"testing"
	"time"

	"../golang"
)

/*
Connection pattern : nodes B,C,D are connected with A.
A request a store on node all other nodes and then searches for it.
*/
func TestFindValue(t *testing.T) {
	A := d7024e.NewKademlia(8500, "5111111400000000000000000000000000000000")
	A.Start()
	time.Sleep(50 * time.Millisecond)

	B := d7024e.NewKademlia(8501, "5111111400000000000000000000000000000001")
	B.Start()
	B.Ping("localhost:8500")
	time.Sleep(50 * time.Millisecond)

	C := d7024e.NewKademlia(8502, "5111111400000000000000000000000000000002")
	C.Start()
	C.Ping("localhost:8500")
	time.Sleep(50 * time.Millisecond)

	D := d7024e.NewKademlia(8503, "5111111400000000000000000000000000000003")
	D.Start()
	D.Ping("localhost:8500")
	time.Sleep(50 * time.Millisecond)

	fmt.Println("All nodes connected")
	node2 := d7024e.NewKademliaID("5111111400000000000000000000000000000000")
	data := []byte("Testing a fucking shit send.")
	A.SendStoreMessage(node2, data)
	time.Sleep(50 * time.Millisecond)

	fmt.Println("Complete store.")
	//After storing an item on node 8401, look it up.
	find := A.SendFindValueMessage(d7024e.NewKademliaID("5111111400000000000000000000000000000000"))
	time.Sleep(50 * time.Millisecond)
	if string(find) == string(""){
		t.Error("Couldn't find the stored value.", find)
	} else {
		fmt.Println("Item returned : ", string(find), "\n")
	}
}
