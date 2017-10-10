package d7024e

import (
	"fmt"
	"sort"
	"sync"
)

type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
}

func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}

func (contact *Contact) Equals(otherContact Contact) bool {
	if contact.ID.Equals(otherContact.ID) && contact.Address == otherContact.Address {
		return true
	} else {
		return false
	}
}

func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.distance.Less(otherContact.distance)
}

func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

func CalcDistances(list *[]Contact, target *KademliaID) {
	for i := 0; i < len(*list); i++ {
		l := *list
		l[i].CalcDistance(target)
	}
}

type ContactCandidates struct {
	contacts []Contact
}

func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	if candidates.Len() >= count {
		return candidates.contacts[:count]
	}
	return candidates.contacts
}

func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(&candidates.contacts[j])
}

type ContactStateItem struct {
	contact  Contact
	queried  bool
	received bool
	counter  int
}

func NewContactStateItem(contact Contact) ContactStateItem {
	return ContactStateItem{contact: contact, queried: false, received: false, counter: 0}
}

func (item *ContactStateItem) Less(contactItem ContactStateItem) bool {
	return item.contact.Less(&contactItem.contact)
}

func (item *ContactStateItem) Equals(contactItem ContactStateItem) bool {
	return item.contact.Equals(contactItem.contact)
}

/*
* A temporary table of contacts.
* Used when needed to organize received contacts asynchronously.
* Is threadsafe.
 */
type ContactStateList struct {
	contacts   []ContactStateItem
	mutex      sync.Mutex
	target     *KademliaID
	k          int
	maxQueries int
	replenishTarget *Contact
}

func (list *ContactStateList) setReplenishTarget(contact Contact){
	//fmt.Println("1")
	list.mutex.Lock()
	defer list.mutex.Unlock()
	if list.replenishTarget == nil {
		//fmt.Println("2")
		list.replenishTarget = &contact
	}else{
		//fmt.Println("Contact in set : ", contact, "\n", list.replenishTarget)
		//if &contact != nil{
		if list.replenishTarget.Less(&contact){
			list.replenishTarget = &contact
		}
	}
	//fmt.Println("3")
}

func (list *ContactStateList) getReplenishTarget() *Contact {
	//fmt.Println("get : ", list.replenishTarget)
	return list.replenishTarget
}


func NewContactStateList(target *KademliaID, k int) ContactStateList {
	return ContactStateList{contacts: []ContactStateItem{}, mutex: sync.Mutex{}, target: target, k: k, maxQueries: 1}
}

/*
* returns k closest contacts to target. if k is larger than list, all contacts will be returned.
* Only returns contacts that responded.
 */
func (list *ContactStateList) GetKClosestContacts() []Contact {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	contacts := []Contact{}
	counter := 0
	for i := 0; i < len(list.contacts); i++ {
		if list.contacts[i].received {
			contacts = append(contacts, list.contacts[i].contact)
			counter += 1
			if counter == list.k {
				return contacts
			}
		}
	}
	return contacts
}

/*
* Marks a contact received. this way we know that is does not need to be queried again.
 */
func (list *ContactStateList) MarkReceived(contact Contact) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	for i := 0; i < len(list.contacts); i++ {
		if list.contacts[i].contact.Equals(contact) {
			list.contacts[i].received = true
			return
		}
	}
}

/*
* Sets a contact not querid, so it can be queried again.
 */
func (list *ContactStateList) SetNotQueried(contact Contact) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	for i := 0; i < len(list.contacts); i++ {
		if list.contacts[i].contact.Equals(contact) {
			list.contacts[i].queried = false
			list.contacts[i].counter += 1
			return
		}
	}
}

/*
* Returns the contact that is next to be queried.
* returns nil if no contact can be queried.
 */
func (list *ContactStateList) GetNextToQuery() *Contact {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	for i := 0; i < len(list.contacts); i++ {
		if list.contacts[i].queried == false && list.contacts[i].counter < list.maxQueries && list.contacts[i].received == false {
			list.contacts[i].queried = true
			return &list.contacts[i].contact
		}
	}
	return nil
}

/*
* returns true if the list is considered finished.
 */
func (list *ContactStateList) Finished() bool {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	counter := 0

	for i := 0; i < len(list.contacts); i++ {
		if list.contacts[i].counter == list.maxQueries {
			continue
		}
		if list.contacts[i].received == false {
			return false
		}
		counter += 1
		if counter == list.k {
			return true
		}
	}
	return true
}

/*
* Appends the contacts, deletes duplicates and sorts them.
 */
func (list *ContactStateList) AppendUniqueSorted(contacts []Contact) {
	list.mutex.Lock()
	CalcDistances(&contacts, list.target)
	for i := 0; i < len(contacts); i++ {
		contact := NewContactStateItem(contacts[i])
		if !list.member(contact) {
			list.sortedInsert(contact)
		}
	}
	list.mutex.Unlock()
}

/*
* returns true if contact is a member of the list.
 */
func (list *ContactStateList) member(contact ContactStateItem) bool {
	for i := 0; i < len(list.contacts); i++ {
		if list.contacts[i].Equals(contact) {
			return true
		}
	}
	return false
}

/*
* inserts contacts into list.contacts, and keeps it sorted.
 */
func (list *ContactStateList) sortedInsert(contact ContactStateItem) {
	l := len(list.contacts)
	if l == 0 {
		list.contacts = []ContactStateItem{contact}
		return
	}
	i := sort.Search(l, func(i int) bool { return contact.Less(list.contacts[i]) }) //list.contacts[i].Less(contact)
	if i == 0 {
		list.contacts = append([]ContactStateItem{contact}, list.contacts...)
		return
	}
	if i == l { // new value is the biggest
		list.contacts = append(list.contacts, contact)
		return
	}
	list.contacts = append(list.contacts[:i], append([]ContactStateItem{contact}, list.contacts[i:]...)...)
}

func (list *ContactStateList) Len() int {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return len(list.contacts)
}
