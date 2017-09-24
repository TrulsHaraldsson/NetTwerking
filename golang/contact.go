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

func Member(list []Contact, contact Contact) bool {
	for i := 0; i < len(list); i++ {
		if list[i].Equals(contact) {
			return true
		}
	}
	return false
}

type ContactCandidates struct {
	contacts []Contact
}

func (canditates *ContactCandidates) Member(contact Contact) bool {
	return Member(canditates.contacts, contact)
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

func (candidates *ContactCandidates) RemoveDuplicates() {
	// Use map to record duplicates as we find them.
	result := []Contact{}

	for v := range candidates.contacts {
		if !Member(result, candidates.contacts[v]) {
			result = append(result, candidates.contacts[v])
			// Do add duplicate.
		}
	}
	// Return the new slice.
	candidates.contacts = result
}

/*
* A temporary table of contacts.
* Used when needed to organize received contacts asynchronously.
* Is threadsafe.
 */
type TempContactTable struct {
	contacts       ContactCandidates
	bannedContacts ContactCandidates
	mutex          sync.Mutex
	target         *KademliaID
	appends        int
}

func NewTempContactTable(target *KademliaID) TempContactTable {
	return TempContactTable{contacts: ContactCandidates{}, bannedContacts: ContactCandidates{}, mutex: sync.Mutex{}, target: target, appends: 0}
}

/*
* Appends the contacts, deletes duplicates and sorts them.
 */
func (table *TempContactTable) AppendUniqueSorted(contacts []Contact) {
	table.mutex.Lock()
	CalcDistances(&contacts, table.target)
	table.contacts.Append(contacts)
	table.contacts.RemoveDuplicates()
	table.contacts.Sort()
	table.appends += 1
	table.mutex.Unlock()
}

func (table *TempContactTable) GetNumOfAppends() int {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	return table.appends
}

func (table *TempContactTable) GetClosest() Contact {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	return table.contacts.contacts[0]
}

func (table *TempContactTable) TargetFound() bool {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	return table.contacts.contacts[0].ID.Equals(table.target)
}

func (table *TempContactTable) Get(index int) Contact {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	return table.contacts.contacts[index]
}

func (table *TempContactTable) GetKContacts(k int) []Contact {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	return table.contacts.GetContacts(k)
}

func (table *TempContactTable) isBanned(contact Contact) bool {
	return table.bannedContacts.Member(contact)
}

func (table *TempContactTable) BannContact(contact Contact) bool {
	table.mutex.Lock()
	if table.isBanned(contact) {
		table.mutex.Unlock()
		return false
	} else {
		table.bannedContacts.Append([]Contact{contact})
		table.mutex.Unlock()
		return true
	}
}

func (table *TempContactTable) Len() int {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	return table.contacts.Len()
}
