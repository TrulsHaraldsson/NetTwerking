package d7024e

import (
	"fmt"
	"testing"
)

func TestRoutingTableBSTUpdate(t *testing.T) {
	me := NewContact(NewKademliaID("F000000000000000000000000000000000000000"), "localhost:8001")
	rt := NewRoutingTableBST(me)

	c1 := NewContact(NewKademliaID("3000000000000000000000000000000000000000"), "localhost:8001")
	rt.Update(c1)

	c2 := NewContact(NewKademliaID("7000000000000000000000000000000000000000"), "localhost:8001")
	rt.Update(c2)

	bucketLeft, nodeLeft := rt.root.findBucket(0, me.ID)
	if rt.root.Left != nodeLeft {
		t.Error("Expected root's left child to exist")
	}

	if bucketLeft.getContact(me.ID) == nil {
		t.Error("Expected to find ourself in this bucket")
	}

	bucketRight, nodeRight := rt.root.findBucket(0, c2.ID)
	if rt.root.Right != nodeRight {
		t.Error("Expected root's right child to exist")
	}

	if bucketRight.getContact(c2.ID) == nil {
		t.Error("Expected to find other contact in this bucket")
	}

	if bucketRight.Len() != 2 {
		t.Error("Expected bucket to have size two")
	}

	if rt.root.Bucket != nil {
		t.Error("Expected root bucket to be nil after split")
	}
}

func TestBSTNewNode(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8001")
	bucket := NewBucket(&me)
	node := NewNode(bucket)

	if node.Left != nil {
		t.Error("Expected Left-Tree to be nil in a newly created Tree")
	}

	if node.Right != nil {
		t.Error("Expected Right-Tree to be nil in newly created Tree")
	}

	if !(node.Bucket.Front().ID.Equals(NewKademliaID("0000000000000000000000000000000000000000"))) {
		t.Error("Expected ID 0000000000000000000000000000000000000000, got", node.Bucket.Front().ID)
	}
}

func TestBSTisLeaf(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8001")
	bucket := NewBucket(&me)
	node := NewNode(bucket)

	if node.isLeaf() == false {
		t.Error("Expected node to be a leaf.")
	}
}

func TestBSTFindClosestContacts(t *testing.T) {
	c1 := NewContact(NewKademliaID("F000000000000000000000000000000000000000"), "localhost:8001")
	c2 := NewContact(NewKademliaID("6000000000000000000000000000000000000000"), "localhost:8001")
	c3 := NewContact(NewKademliaID("5000000000000000000000000000000000000000"), "localhost:8001")
	c4 := NewContact(NewKademliaID("4000000000000000000000000000000000000000"), "localhost:8001")
	c5 := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8001")

	me := NewContact(NewKademliaID("FF00000000000000000000000000000000000000"), "localhost:8001")

	b1 := NewBucket(&c1)
	b2 := NewBucket(&c2)
	b3 := NewBucket(&c3)
	b4 := NewBucket(&c4)
	b5 := NewBucket(&c5)

	rt := NewRoutingTableBST(me)
	rt.root.Bucket = nil

	rootL := NewNode(b1)
	rootL.Parent = rt.root
	rootR := NewNode(nil)
	rootR.Parent = rt.root
	rt.root.Left = rootL
	rt.root.Right = rootR

	rootRR := NewNode(b5)
	rootRR.Parent = rootR
	rootRL := NewNode(nil)
	rootRL.Parent = rootR
	rootR.Left = rootRL
	rootR.Right = rootRR

	rootRLL := NewNode(b2)
	rootRLL.Parent = rootRL
	rootRLR := NewNode(nil)
	rootRLR.Parent = rootRL
	rootRL.Left = rootRLL
	rootRL.Right = rootRLR

	rootRLRL := NewNode(b3)
	rootRLRL.Parent = rootRLR
	rootRLRR := NewNode(b4)
	rootRLRR.Parent = rootRLR
	rootRLR.Left = rootRLRL
	rootRLR.Right = rootRLRR

	contacts := rt.FindClosestContacts(c3.ID, 3)
	if len(contacts) != 3 {
		t.Error("Expected the length of returned array to be 3, got", len(contacts))
	}

	contacts = rt.FindClosestContacts(c3.ID, 20)
	if len(contacts) != 5 {
		t.Error("Expected the length of returned array to be 5, got", len(contacts))
	}
}

func TestBSTNext(t *testing.T) {
	c1 := NewContact(NewKademliaID("F000000000000000000000000000000000000000"), "localhost:8001")
	c2 := NewContact(NewKademliaID("6000000000000000000000000000000000000000"), "localhost:8001")
	c3 := NewContact(NewKademliaID("5000000000000000000000000000000000000000"), "localhost:8001")
	c4 := NewContact(NewKademliaID("4000000000000000000000000000000000000000"), "localhost:8001")
	c5 := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8001")

	b1 := NewBucket(&c1)
	b2 := NewBucket(&c2)
	b3 := NewBucket(&c3)
	b4 := NewBucket(&c4)
	b5 := NewBucket(&c5)

	root := NewNode(nil)

	rootL := NewNode(b1)
	rootL.Parent = root
	rootR := NewNode(nil)
	rootR.Parent = root
	root.Left = rootL
	root.Right = rootR

	rootRR := NewNode(b5)
	rootRR.Parent = rootR
	rootRL := NewNode(nil)
	rootRL.Parent = rootR
	rootR.Left = rootRL
	rootR.Right = rootRR

	rootRLL := NewNode(b2)
	rootRLL.Parent = rootRL
	rootRLR := NewNode(nil)
	rootRLR.Parent = rootRL
	rootRL.Left = rootRLL
	rootRL.Right = rootRLR

	rootRLRL := NewNode(b3)
	rootRLRL.Parent = rootRLR
	rootRLRR := NewNode(b4)
	rootRLRR.Parent = rootRLR
	rootRLR.Left = rootRLRL
	rootRLR.Right = rootRLRR

	if rootL.next() != rootRLL {
		t.Error("Expected rootRLL to be next")
	}

	if rootRLL.next() != rootRLRL {
		t.Error("Expected rootRLRL to be next")
	}

	if rootRLRL.next() != rootRLRR {
		t.Error("Expected rootRLRR to be next")
	}

	if rootRLRR.next() != rootRR {
		t.Error("Expected rootRR to be next")
	}

	if rootRR.next() != nil {
		t.Error("Expected rootRR's next to be nil")
	}
}

func TestBSTPrev(t *testing.T) {
	c1 := NewContact(NewKademliaID("F000000000000000000000000000000000000000"), "localhost:8001")
	c2 := NewContact(NewKademliaID("6000000000000000000000000000000000000000"), "localhost:8001")
	c3 := NewContact(NewKademliaID("5000000000000000000000000000000000000000"), "localhost:8001")
	c4 := NewContact(NewKademliaID("4000000000000000000000000000000000000000"), "localhost:8001")
	c5 := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8001")

	b1 := NewBucket(&c1)
	b2 := NewBucket(&c2)
	b3 := NewBucket(&c3)
	b4 := NewBucket(&c4)
	b5 := NewBucket(&c5)

	root := NewNode(nil)

	rootL := NewNode(b1)
	rootL.Parent = root
	rootR := NewNode(nil)
	rootR.Parent = root
	root.Left = rootL
	root.Right = rootR

	rootRR := NewNode(b5)
	rootRR.Parent = rootR
	rootRL := NewNode(nil)
	rootRL.Parent = rootR
	rootR.Left = rootRL
	rootR.Right = rootRR

	rootRLL := NewNode(b2)
	rootRLL.Parent = rootRL
	rootRLR := NewNode(nil)
	rootRLR.Parent = rootRL
	rootRL.Left = rootRLL
	rootRL.Right = rootRLR

	rootRLRL := NewNode(b3)
	rootRLRL.Parent = rootRLR
	rootRLRR := NewNode(b4)
	rootRLRR.Parent = rootRLR
	rootRLR.Left = rootRLRL
	rootRLR.Right = rootRLRR

	if rootRLL.prev() != rootL {
		t.Error("Expected rootL to be prev")
	}

	if rootRLRL.prev() != rootRLL {
		t.Error("Expected rootRLL to be prev")
	}

	if rootRLRR.prev() != rootRLRL {
		t.Error("Expected rootRLRL to be prev")
	}

	if rootRR.prev() != rootRLRR {
		t.Error("Expected rootRLRR to be prev")
	}

	if rootL.prev() != nil {
		t.Error("Expected rootL's prev to be nil")
	}
}

func TestBSTFindBucket(t *testing.T) {
	c1 := NewContact(NewKademliaID("F000000000000000000000000000000000000000"), "localhost:8001")
	c2 := NewContact(NewKademliaID("6000000000000000000000000000000000000000"), "localhost:8001")
	c3 := NewContact(NewKademliaID("5000000000000000000000000000000000000000"), "localhost:8001")
	c4 := NewContact(NewKademliaID("4000000000000000000000000000000000000000"), "localhost:8001")
	c5 := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8001")

	b1 := NewBucket(&c1)
	b2 := NewBucket(&c2)
	b3 := NewBucket(&c3)
	b4 := NewBucket(&c4)
	b5 := NewBucket(&c5)

	root := NewNode(nil)

	rootL := NewNode(b1)
	rootL.Parent = root
	rootR := NewNode(nil)
	rootR.Parent = root
	root.Left = rootL
	root.Right = rootR

	rootRR := NewNode(b5)
	rootRR.Parent = rootR
	rootRL := NewNode(nil)
	rootRL.Parent = rootR
	rootR.Left = rootRL
	rootR.Right = rootRR

	rootRLL := NewNode(b2)
	rootRLL.Parent = rootRL
	rootRLR := NewNode(nil)
	rootRLR.Parent = rootRL
	rootRL.Left = rootRLL
	rootRL.Right = rootRLR

	rootRLRL := NewNode(b3)
	rootRLRL.Parent = rootRLR
	rootRLRR := NewNode(b4)
	rootRLRR.Parent = rootRLR
	rootRLR.Left = rootRLRL
	rootRLR.Right = rootRLRR

	bucket, node := root.findBucket(0, c1.ID)
	if node != rootL {
		t.Error("Expected rootL to be found")
	}

	if bucket.getContact(c1.ID) == nil {
		t.Error("Expected to find c1 in this bucket")
	}

	bucket, node = root.findBucket(0, c4.ID)
	if node != rootRLRR {
		t.Error("Expected rootRLRR to be found")
	}

	if bucket.getContact(c4.ID) == nil {
		t.Error("Expected to find c4 in this bucket")
	}
}

func TestMyBucket(t *testing.T) {
	c1 := NewContact(NewKademliaID("F000000000000000000000000000000000000000"), "localhost:8001")
	c2 := NewContact(NewKademliaID("7000000000000000000000000000000000000000"), "localhost:8001")

	bucket := NewBucket(&c1)

	if bucket.getContact(c1.ID) == nil {
		t.Error("Expected to find contact c1 in bucket, got nil")
	}

	if bucket.getContact(c2.ID) != nil {
		t.Error("Did not expect to find contact c2 in bucket")
	}
}

func TestGetNBit(t *testing.T) {
	id := NewKademliaID("800000000000000000000000000000000000000E") // 1000 ... 1110

	res := getNBit(0, id)
	if res != 1 {
		t.Error("Expected bit 0 to have value 1, got ", res)
	}

	res = getNBit(1, id)
	if res != 0 {
		t.Error("Expected bit 1 to have value 0, got ", res)
	}

	res = getNBit(159, id)
	if res != 0 {
		t.Error("Expected bit 0 to have value 1, got ", res)
	}

	res = getNBit(158, id)
	if res != 1 {
		t.Error("Expected bit 0 to have value 1, got ", res)
	}

	res = getNBit(160, id)
	if res != -1 {
		t.Error("Expected return to be -1, got", res)
	}
}

func TestIsNBitsEqual(t *testing.T) {
	fmt.Println("RUNNING")
	id1 := NewKademliaID("5000000000000000000000000000000000000000") // First bits 0101 ...
	id2 := NewKademliaID("6000000000000000000000000000000000000000") // First bits 0110 ...

	if isNBitsEqual(1, id1, id2) == false {
		t.Error("Expected the first bit to be equal")
	}

	if isNBitsEqual(2, id1, id2) == false {
		t.Error("Expected the first two bits to be equal")
	}

	if isNBitsEqual(3, id1, id2) == true {
		t.Error("Expected third bit to be not equal")
	}

	if isNBitsEqual(4, id1, id2) == true {
		t.Error("Expected the forth bit to be not equal")
	}
}

func TestBSTOneNode(t *testing.T) {
	root := NewNode(nil)
	root.prev()
	root.next()

}
