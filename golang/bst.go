package d7024e

import (
	"container/list"
	"sync"
)

/*
 * A routing table have a Contact 'me' which is the first entry in the table.
 * The routing table also holds a reference to the root of a binary search tree.
 */
type RoutingTable struct {
	me   *Contact
	root *Node
	mux  sync.Mutex
}

/*
 * Creates a new routingtable and returns a pointer to it.
 * The given Contact will be the initial entry in the table.
 */
func newRoutingTable(me Contact) *RoutingTable {
	rt := RoutingTable{}
	rt.me = &me
	bucket := newBucket(&me)
	rt.root = newNode(bucket)
	return &rt
}

/*
 * Given a contact it either inserts the contact in appropriate bucket, or
 * if the contact allready exists in the routing table the contact will be
 * moved to the front of the bucket. Lastly, if the appropriate bucket is full
 * the contact will simply be discarded.
 */
func (this *RoutingTable) update(contact Contact) {
	this.mux.Lock()
	defer this.mux.Unlock()
	bucket, node := this.root.findBucket(0, contact.ID)

	if bucket.front().Equals(*this.me) {
		for i := 0; i < 160; i++ {
			meBit := getNBit(uint(i), this.me.ID)
			otherBit := getNBit(uint(i), contact.ID)
			if meBit != otherBit {
				if meBit == 1 {
					left := newNode(node.Bucket)
					left.Parent = node
					node.Left = left
					node.Bucket = nil
					right := newNode(newBucket(&contact))
					right.Parent = node
					node.Right = right
				} else {
					left := newNode(newBucket(&contact))
					left.Parent = node
					node.Left = left
					right := newNode(node.Bucket)
					right.Parent = node
					node.Right = right
					node.Bucket = nil
				}
				return
			}

		}
	} else {
		if bucket.Len() < 20 {
			c := bucket.getContact(contact.ID)
			if c == nil {
				bucket.list.PushFront(contact)
			} else {
				bucket.list.MoveToFront(c)
			}
		}
	}
}

func (this *RoutingTable) findClosestContacts(target *KademliaID, count int) []Contact {
	this.mux.Lock()
	defer this.mux.Unlock()
	var candidates ContactCandidates
	bucket, node := this.root.findBucket(0, target)

	candidates.Append(bucket.getContactAndCalcDistance(target))

	prev := node.prev()
	next := node.next()

	for {
		if candidates.Len() >= count {
			break
		} else {
			if prev == nil && next == nil {
				break
			}

			if next != nil {
				candidates.Append(next.Bucket.getContactAndCalcDistance(target))
				next = next.next()
			}

			if prev != nil {
				candidates.Append(prev.Bucket.getContactAndCalcDistance(target))
				prev = prev.prev()
			}
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
}

/*
 * A node in the Binary Search Tree. Contain a single bucket and
 * two pointers to its children.
 */
type Node struct {
	Left   *Node
	Right  *Node
	Parent *Node
	Bucket *MyBucket
}

/*
 * Creates a new node and sets it bucket to the given bucket
 */
func newNode(bucket *MyBucket) *Node {
	node := Node{}
	node.Bucket = bucket
	node.Left = nil
	node.Right = nil
	node.Parent = nil
	return &node
}

/*
 * Search for a bucket that covers the given contact's ID. The node that
 * holds the bucket is also returned.
 */
func (this *Node) findBucket(index int, kademliaID *KademliaID) (*MyBucket, *Node) {
	if this.isLeaf() {
		return this.Bucket, this
	} else {
		// Compare i'th bit
		bit := getNBit(uint(index), kademliaID)
		if bit == 0 {
			return this.Right.findBucket(index+1, kademliaID)
		} else {
			return this.Left.findBucket(index+1, kademliaID)
		}
	}
}

/*
 * Returns
 */
func (this *Node) next() *Node {
	if this.isLeftChild() {
		next := this.Parent.Right
		for {
			if next.isLeaf() {
				return next
			} else {
				next = next.Left
			}
		}
	} else {
		next := this.Parent
		if next == nil {
			return nil
		}
		for {
			if next.Parent == nil {
				return nil
			}

			if next.isRightChild() {
				next = next.Parent
			} else {
				return next.Parent.Right
			}
		}
	}
}

/*
 * Return the
 */
func (this *Node) prev() *Node {
	if this.isRightChild() {
		prev := this.Parent.Left
		for {
			if prev.isLeaf() {
				return prev
			} else {
				prev = prev.Right
			}
		}
	} else {
		prev := this.Parent
		if prev == nil {
			return nil
		}
		for {
			if prev.Parent == nil {
				return nil
			}

			if prev.isLeftChild() {
				prev = prev.Parent
			} else {
				return prev.Parent.Left
			}
		}
	}
}

/*
 * Returns True if this node is the left child of it's parent, else False.
 */
func (this *Node) isLeftChild() bool {
	if this.Parent != nil && this.Parent.Left == this {
		return true
	} else {
		return false
	}
}

/*
 * Return True if this node is the right child of it's parent, else False.
 */
func (this *Node) isRightChild() bool {
	if this.Parent != nil && this.Parent.Right == this {
		return true
	} else {
		return false
	}
}

/*
 * A Node is a leaf if both of its child nodes are nil.
 */
func (this *Node) isLeaf() bool {
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
func newBucket(contact *Contact) *MyBucket {
	bucket := MyBucket{list.New()}
	bucket.list.PushFront(*contact)
	return &bucket
}

func (this *MyBucket) getContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := this.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

/*
 * Returns the contact that have been the most recently seen.
 */
func (this *MyBucket) front() *Contact {
	contact := this.list.Front().Value.(Contact)
	return &contact
}

/*
 * Returns the contact with the given ID if it is in the bucket. If
 * no contact with that ID exists in the bucket the return will be nil.
 */
func (this *MyBucket) getContact(id *KademliaID) *list.Element {
	for e := this.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID
		if id.Equals(nodeID) {
			return e
		}
	}
	return nil
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
func getNBit(n uint, id *KademliaID) int {
	if 0 <= n && n <= 159 {
		byteToCheck := id[n/8] // 8 bits per byte
		bitToCheck := n % 8    // 8 Bits
		res := byteToCheck & (1 << (7 - bitToCheck))
		//fmt.Println("ByteToCheck", byteToCheck)
		//fmt.Println("bitToCheck", bitToCheck)
		//fmt.Println("Res:", res)
		if res > 0 {
			return 1
		} else {
			return 0
		}
	} else {
		return -1
	}
}

/*
 * Compares the two given KademliaIDs and Returns TRUE if
 * the first N bits are equal, else FALSE
 */
func isNBitsEqual(n uint, id1 *KademliaID, id2 *KademliaID) bool {
	var i uint
	for i = 0; i < n; i++ {
		if getNBit(i, id1) != getNBit(i, id2) {
			return false
		}
	}
	return true
}
