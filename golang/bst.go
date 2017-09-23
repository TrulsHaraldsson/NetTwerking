package d7024e

import (
	"container/list"
	"errors"
	"fmt"
)

/*
 * A routing table have a Contact 'me' which is the first entry in the table. 
 * The routing table also holds a reference to the root of a binary search tree.
 */
type RoutingTableBST struct {
	me *Contact
	root *Node
}

/*
 * Creates a new routingtable and returns a pointer to it. 
 * The given Contact will be the initial entry in the table.
 */ 
func NewRoutingTableBST(me Contact) *RoutingTableBST {
	rt := RoutingTableBST{}
	rt.me = &me
	bucket := NewBucket(&me)
	rt.root = NewNode(bucket)
	return &rt
}

/*
 * Given a contact it either inserts the contact in appropriate bucket, or
 * if the contact allready exists in the routing table the contact will be 
 * moved to the front of the bucket. Lastly, if the appropriate bucket is full
 * the contact will simply be discarded.
 */
func (this *RoutingTableBST) Update(contact Contact) {

}

/*
 * A node in the Binary Search Tree. Contain a single bucket and
 * two pointers to its children.
 */
type Node struct {
	Left *Node
	Right *Node
	Bucket *MyBucket
}

/*
 * Creates a new node and sets it bucket to the given bucket
 */
func NewNode(bucket *MyBucket) *Node{
	node := Node{}
	node.Bucket = bucket
	node.Left = nil
	node.Right = nil
	return &node
}

func (this *Node) insert(contact *Contact) {
	//id := contact.ID
	//bit := getNBit(contact.ID)
	
	// Search for Right bucket
	// When found Insert
	// If not found, Create new?
	// If original 
}

/*
 * A Node is a leaf if both of its child nodes are nil.
 */
func (this *Node) IsLeaf() bool {
	return this.Left == nil && this.Right == nil
}


/*
 * A bucket holds contact information. One bucket can hold at 
 * maximum 20 contacts.  
 */
type MyBucket struct {
	list *list.List
}

/*
 * Creates a new bucket and inserts the given contact.
 * Returns the pointer to the bucket.
 */
func NewBucket(contact *Contact) *MyBucket {
	bucket := MyBucket{list.New()}
	bucket.list.PushFront(*contact)
	return &bucket
}

/*
 * Returns the contact that have been the most recently seen.
 */
func (this *MyBucket) Front() *Contact {
	contact := this.list.Front().Value.(Contact)
	return &contact
}

/*
 * Returns the number of KademliaIDs stored in this bucket.
 */
func (this *MyBucket) Len() int {
	return this.list.Len()
}

/*
 * Returns the value of the specified bit, either 0 or 1. If the specified
 * bit is not between 0 and 159 the result will be -1. 
 */
func GetNBit(n uint, id *KademliaID) (int, error) {
	if 0 <= n && n <= 159 {
		byteToCheck := id[n / 8]; // 8 bits per byte
		bitToCheck := n % 8 // 8 Bits
		res := byteToCheck & (1<<(7-bitToCheck))
		fmt.Println("ByteToCheck", byteToCheck)
		fmt.Println("bitToCheck", bitToCheck)
		fmt.Println("Res:", res)
		if res > 0 {
			return 1, nil
		}else {
			return 0, nil
		}
	}else {
		return -1, errors.New("Input has to be from 0 and up to 159")
	}
}
