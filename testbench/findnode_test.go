package main

import (
	"fmt"
	"testing"
	"time"

	"../golang"
)

func TestFindNode(t *testing.T) {
	d7024e.StartNode(8100, "none", "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	time.Sleep(10 * time.Millisecond)                                                          //B
	n1 := d7024e.StartNode(8101, "localhost:8100", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA") //A
	d7024e.StartNode(8102, "localhost:8100", "CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")       //C
	time.Sleep(10 * time.Millisecond)
	d7024e.StartNode(8103, "localhost:8102", "DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD") //D
	fmt.Println("All nodes connected", n1)
	contact := n1.SendFindContactMessage(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD"))

	if !contact.ID.Equals(d7024e.NewKademliaID("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")) {
		t.Error("Not correct contact returned", contact)
	}
}
