package main

import (
	"testing"
	"time"

	"../golang"
)

/*
Connection pattern : nodes B,C,D are connected with A.
A request a store on node all other nodes and then searches for it.
*/
func TestFindValue(t *testing.T) {
	A := d7024e.NewKademlia("localhost:8500", "5111111400000000000000000000000000000000")
	A.StartListening()
	time.Sleep(50 * time.Millisecond)

	B := d7024e.NewKademlia("localhost:8501", "5111111400000000000000000000000000000001")
	B.StartListening()
	B.Ping("localhost:8500")
	time.Sleep(50 * time.Millisecond)

	C := d7024e.NewKademlia("localhost:8502", "5111111400000000000000000000000000000002")
	C.StartListening()
	C.Ping("localhost:8500")
	time.Sleep(50 * time.Millisecond)

	D := d7024e.NewKademlia("localhost:8503", "5111111400000000000000000000000000000003")
	D.StartListening()
	D.Ping("localhost:8500")
	time.Sleep(50 * time.Millisecond)

	//Create file first
	filename2 := "filename2"
	fileKID := d7024e.NewValueID(&filename2)
	fileID := fileKID.String()
	data2 := []byte("Testing a send.")
	A.Store(fileKID, &data2)
	time.Sleep(50 * time.Millisecond)

	strA := A.SearchFileLocal(&fileID)
	if *strA != string(data2) {
		t.Error("Strings of content dont match!")
	}
	strB := B.SearchFileLocal(&fileID)
	if *strB != string(data2) {
		t.Error("Strings of content dont match!")
	}
	strC := C.SearchFileLocal(&fileID)
	if *strC != string(data2) {
		t.Error("Strings of content dont match!")
	}
	strD := D.SearchFileLocal(&fileID)
	if *strD != string(data2) {
		t.Error("Strings of content dont match!")
	}

	E := d7024e.NewKademlia("localhost:8504", "5111111700000000000000000000000000000004")
	E.StartListening()
	E.Ping("localhost:8500")
	time.Sleep(50 * time.Millisecond)
	find := E.FindValue(&filename2)
	if find == nil {
		t.Error("Not found!")
	}
	contentReceived := string(find)
	if contentReceived != string(data2) {
		t.Error("Strings of content dont match!")
	}
	time.Sleep(time.Millisecond * 20)
	E.DeleteFileLocal(filename2)
}
