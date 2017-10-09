package main

import (
	"testing"
	"time"

	"../golang"
)

func TestBootstrap1(t *testing.T) {
	//"7a9eb1929b4615f8886a229ae273c649d7c4d3ab"
	//"6bc2159410a7e14865e82baa0accaa958801e613"
	addr := "localhost:8510"
	kID := "6bc2159410a7e14865e82baa0accaa958801e613"
	c1 := d7024e.NewContact(d7024e.NewKademliaID(kID), addr)
	k1 := d7024e.CreateAndStartNode(addr, kID, nil)
	time.Sleep(50 * time.Millisecond)

	k2 := d7024e.CreateAndStartNode("localhost:8511", "none", &c1)
	time.Sleep(50 * time.Millisecond)
	if k2.RT.Contacts() != 2 {
		t.Error("Wrong amount of contacts: ", k2.RT.Contacts())
	}
	k3 := d7024e.CreateAndStartNode("localhost:8512", "none", &c1)
	time.Sleep(50 * time.Millisecond)
	if k3.RT.Contacts() != 3 {
		t.Error("Wrong amount of contacts: ", k3.RT.Contacts())
	}

	if k1.RT.Contacts() != 3 {
		t.Error("Wrong amount of contacts: ", k1.RT.Contacts())
	}

}

func TestBootstrap2(t *testing.T) {
	addr := "localhost:8610"
	kID := "6bc2159410a7e14865e82baa0accaa958801e613"
	c1 := d7024e.NewContact(d7024e.NewKademliaID(kID), addr)
	d7024e.CreateAndStartNode(addr, kID, nil)
	time.Sleep(50 * time.Millisecond)

	count := 50
	nodes := make([]d7024e.Kademlia, count)
	for i := 0; i < count; i++ {
		addr := d7024e.CreateAddr("localhost", 8611+i)
		k := d7024e.CreateAndStartNode(addr, "none", &c1)
		nodes[i] = *k
		//time.Sleep(50 * time.Millisecond)
	}
	for _, val := range nodes {
		if val.RT.Contacts() < 20 {
			t.Error("amount of contacts to small:", val.RT.Contacts())
		}
	}

}
