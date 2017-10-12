package main

import (
	//"fmt"
	"testing"
	"time"

	"../golang"
)

/*
Connection pattern : nodes B,C,D are connected with A.
A requests a store on all nodes.
*/
func TestStoreToAll(t *testing.T) {
	A := d7024e.NewKademlia("localhost:8400", "2111111400000000000000000000000000000000")
	A.StartListening()
	time.Sleep(50 * time.Millisecond)

	B := d7024e.NewKademlia("localhost:8401", "2111111400000000000000000000000000000001")
	B.StartListening()
	B.Ping("localhost:8400")
	time.Sleep(50 * time.Millisecond)

	C := d7024e.NewKademlia("localhost:8402", "2111111400000000000000000000000000000002")
	C.StartListening()
	C.Ping("localhost:8400")
	time.Sleep(50 * time.Millisecond)

	D := d7024e.NewKademlia("localhost:8403", "2111111400000000000000000000000000000003")
	D.StartListening()
	D.Ping("localhost:8400")
	time.Sleep(50 * time.Millisecond)

	//fmt.Println("All nodes connected")
	filename := "filename3"
	fileKID := d7024e.NewValueID(&filename)
	fileID := fileKID.String()
	data := []byte("Testing a send3.")
	A.Store(fileKID, &data)
	time.Sleep(50 * time.Millisecond)
	file := D.SearchFileLocal(&fileID)
	if string(*file) != string(data) {
		t.Error("Wrong file content")
	}
	A.DeleteFileLocal(fileID)
}

/*
Connection pattern : nodes A and B are connected.
A Sends multiple stores to B.
*/
func TestStoreToOne(t *testing.T) {
	A := d7024e.NewKademlia("localhost:8410", "2111111400000000000000000000000000000010")
	A.StartListening()
	time.Sleep(50 * time.Millisecond)

	B := d7024e.NewKademlia("localhost:8411", "2111111400000000000000000000000000000011")
	B.StartListening()
	B.Ping("localhost:8410")
	time.Sleep(50 * time.Millisecond)

	//fmt.Println("All nodes connected")
	filename1 := "failname"
	fileKID := d7024e.NewValueID(&filename1)
	fileID := fileKID.String()
	data1 := []byte("Testing a send4.")
	A.Store(fileKID, &data1)
	time.Sleep(50 * time.Millisecond)
	file := B.SearchFileLocal(&fileID)
	if string(*file) != string(data1) {
		t.Error("Wrong file content")
	}
	A.DeleteFileLocal(fileID)
}
