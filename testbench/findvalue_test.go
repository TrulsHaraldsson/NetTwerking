package main

import (
	//"fmt"
	//"strings"
	"testing"
	"time"
	"io/ioutil"
	//"reflect"
	"../golang"

)

/*
Connection pattern : nodes B,C,D are connected with A.
A request a store on node all other nodes and then searches for it.
*/
func TestFindValue(t *testing.T) {
	A := d7024e.NewKademlia("localhost:8500", "5111111400000000000000000000000000000000")
	A.Start()
	time.Sleep(50 * time.Millisecond)

	B := d7024e.NewKademlia("localhost:8501", "5111111400000000000000000000000000000001")
	B.Start()
	B.Ping("localhost:8500")
	time.Sleep(50 * time.Millisecond)

	C := d7024e.NewKademlia("localhost:8502", "5111111400000000000000000000000000000002")
	C.Start()
	C.Ping("localhost:8500")
	time.Sleep(50 * time.Millisecond)

	D := d7024e.NewKademlia("localhost:8503", "5111111400000000000000000000000000000003")
	D.Start()
	D.Ping("localhost:8500")
	time.Sleep(50 * time.Millisecond)

	//Create file first
	filename2 := "filename2"
	data2 := []byte("Testing a fucking shit send.")
	filename3 := "./../newfiles/" + string(filename2)
	err2 := ioutil.WriteFile(filename3, text, 0644)
	if err2 != nil {
		panic(err2)
	}

	//Read file
	searchName := "./../newfiles/filename2"
	content, err2 := ioutil.ReadFile(searchName)
	if err2 != nil {
		panic(err2)
	}
	A.SendStoreMessage(&filename2, &content)
	time.Sleep(50 * time.Millisecond)

	find := A.SendFindValueMessage(&filename2)
	if find == nil {
		t.Error("Not found!")
	}
	time.Sleep(50 * time.Millisecond)
}
