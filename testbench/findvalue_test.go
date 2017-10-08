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

	////fmt.Println("All nodes connected")
	filename2 := "filename2"
	//data := []byte("Testing a fucking shit send.")

	searchName := "/tmp/filename2"

	content, err2 := ioutil.ReadFile(searchName)
	if err2 != nil {
		//fmt.Println("ReadMemory if file exist: ", err2)
	}
	//fmt.Println("Loaded complete : ", content)
	A.SendStoreMessage(&filename2, &content)
	time.Sleep(50 * time.Millisecond)

//	//fmt.Println("Complete store!")
	//After storing an item on node 8401, look it up.
	A.SendStoreMessage(&filename2, &content)
	//find := A.SendFindValueMessage(d7024e.NewValueID(&filename2))
	find := A.SendFindValueMessage(&filename2)
	//DO SO YOU CAN SEE CONTENT OF THIS!
	//fmt.Println("Returned File 1 : ", find, "type : ", reflect.TypeOf(find))
	if find == nil {
		t.Error("Not found!")
	}else{
		//file := string(find)
		//fmt.Println("Returned File 2 : ", string(find), "type : ", reflect.TypeOf(find))
		//fmt.Println("Returned File 3 : ", string(file), "type : ", reflect.TypeOf(file))

/*		//fmt.Println(string(find))
		info := string(find)
		r1 := strings.Split(info, ",")
		//fmt.Println(r1[1])
		r1 = strings.Split(r1[1], ":")
		//fmt.Println(r1[1])
		s := strings.TrimRight(r1[1], "}")
		//fmt.Println("Data  : ", s, " type ", reflect.TypeOf(s))
		if len(s) > 0 && s[0] == '"' {
    	s = s[1:]
		}
		if len(s) > 0 && s[len(s)-1] == '"' {
			s = s[:len(s)-1]
		}
		//fmt.Println(s, "type : ", reflect.TypeOf(s))
*/
/*		d := []uint8(s)
		//fmt.Println(d, "type", reflect.TypeOf(d))

		//fmt.Printf("The value of newStr is %s \n", d)
		//fmt.Printf("The type of newStr is %v \n", reflect.TypeOf(d))

		// convert back to string
		backToStr := string([]byte(d[:]))

		//fmt.Println(backToStr)

		//fmt.Printf("The value of backToStr is [%s] \n", backToStr)
		//fmt.Printf("The type of backToStr is %v \n", reflect.TypeOf(backToStr))
*/	}
	time.Sleep(50 * time.Millisecond)
}
