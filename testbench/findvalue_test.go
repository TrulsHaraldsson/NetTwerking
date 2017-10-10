package main

import (
	"io/ioutil"
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
	data2 := []byte("Testing a fucking shit send.")
	filename3 := "../newfiles/" + string(filename2)
	err2 := ioutil.WriteFile(filename3, data2, 0644)
	if err2 != nil {
		panic(err2)
	}
	//Read file
	content, err2 := ioutil.ReadFile(filename3)
	if err2 != nil {
		t.Error("Error when reading!")
	}
	A.Store(&filename2, &content)
	time.Sleep(50 * time.Millisecond)

	strA := A.SearchFileLocal(&filename2)
	if *strA != string(content) {
		t.Error("Strings of content dont match!")
	}
	strB := B.SearchFileLocal(&filename2)
	if *strB != string(content) {
		t.Error("Strings of content dont match!")
	}
	strC := C.SearchFileLocal(&filename2)
	if *strC != string(content) {
		t.Error("Strings of content dont match!")
	}
	strD := D.SearchFileLocal(&filename2)
	if *strD != string(content) {
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
	if contentReceived != string(content) {
		t.Error("Strings of content dont match!")
	}
}
