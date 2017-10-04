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
func TestStoreToAll(t *testing.T) {
	A := d7024e.NewKademlia("localhost", 8400, "2111111400000000000000000000000000000000")
	A.Start()
	time.Sleep(50 * time.Millisecond)

	B := d7024e.NewKademlia("localhost", 8401, "2111111400000000000000000000000000000001")
	B.Start()
	B.Ping("localhost:8400")
	time.Sleep(50 * time.Millisecond)

	C := d7024e.NewKademlia("localhost", 8402, "2111111400000000000000000000000000000002")
	C.Start()
	C.Ping("localhost:8400")
	time.Sleep(50 * time.Millisecond)

	D := d7024e.NewKademlia("localhost", 8403, "2111111400000000000000000000000000000003")
	D.Start()
	D.Ping("localhost:8400")
	time.Sleep(50 * time.Millisecond)

	fmt.Println("All nodes connected")
	filename := "filenameShit"
	data := []byte("Testing a fucking shit send.")
	A.SendStoreMessage(&filename, &data)
	time.Sleep(50 * time.Millisecond)
}

/*
Connection pattern : nodes A and B are connected.
A Sends multiple stores to B.
*/
func TestStoreToOne(t *testing.T) {
	A := d7024e.NewKademlia("localhost", 8410, "2111111400000000000000000000000000000010")
	A.Start()
	time.Sleep(50 * time.Millisecond)

	B := d7024e.NewKademlia("localhost", 8411, "2111111400000000000000000000000000000011")
	B.Start()
	B.Ping("localhost:8410")
	time.Sleep(50 * time.Millisecond)

	fmt.Println("All nodes connected")
	filename1 := "failname"
	data1 := []byte("Testing a fucking shit send 1 time.")
	A.SendStoreMessage(&filename1, &data1)
	time.Sleep(50 * time.Millisecond)
}
