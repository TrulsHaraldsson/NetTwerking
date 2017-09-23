package d7024e

import (
	"testing"
	//"fmt"
)

func TestNodeNewNode (t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8001")
	bucket := NewBucket(&me)
	node := NewNode(bucket)
	
	if node.Left != nil {
		t.Error("Expected Left-Tree to be nil in a newly created Tree")
	}

	if node.Right != nil {
		t.Error("Expected Right-Tree to be nil in newly created Tree")
	}

	if !( node.Bucket.Front().ID.Equals(NewKademliaID("0000000000000000000000000000000000000000")) ) {
		t.Error("Expected ID 0000000000000000000000000000000000000000, got", node.Bucket.Front().ID)
	}
}

func TestNodeIsLeaf(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8001")
	bucket := NewBucket(&me)
	node := NewNode(bucket)

	if node.IsLeaf() == false {
		t.Error("Expected node to be a leaf.")
	}
}

func TestGetNBit(t *testing.T) {
	id := NewKademliaID("800000000000000000000000000000000000000E")
	
	res, _ := GetNBit(0, id)
	if res != 1 {
		t.Error("Expected bit 0 to have value 1, got ", res)
	}

	res, _ = GetNBit(159, id)
	if res != 0 {
		t.Error("Expected bit 0 to have value 1, got ", res)
	}

	res, _ = GetNBit(158, id)
	if res != 1 {
		t.Error("Expected bit 0 to have value 1, got ", res)
	}

	res, _ = GetNBit(160, id)
	if res != -1 {
		t.Error("Expected return to be -1, got", res)
	}
	
}
