package d7024e

import (
	"testing"
)

func TestRoutingTable(t *testing.T) {
	_, rt := CreateTestRT()

	rt.findClosestContacts(
		NewKademliaID("2111111400000000000000000000000000000000"), 20)

}

func CreateTestRT9() ([]Contact, *RoutingTable) {
	rt := newRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))

	for _, contact := range contacts {
		rt.update(contact)
	}
	return contacts, rt
}

func CreateTestRT8() ([]Contact, *RoutingTable) {
	rt := newRoutingTable(NewContact(
		NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003"))
	contacts = append(contacts, NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	for _, contact := range contacts {
		rt.update(contact)
	}
	return contacts, rt
}

func CreateTestRT19() ([]Contact, *RoutingTable) {
	rt := newRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:9102"))

	for _, contact := range contacts {
		rt.update(contact)
	}
	return contacts, rt
}

func CreateTestRT18() ([]Contact, *RoutingTable) {
	rt := newRoutingTable(NewContact(
		NewKademliaID("1111111100000000000000000000000000000000"), "localhost:9102"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	for _, contact := range contacts {
		rt.update(contact)
	}
	return contacts, rt
}

func TestRoutingTableupdate(t *testing.T) {
	rt := newRoutingTable(NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	contacts = append(contacts, NewContact(
		NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))

	for _, contact := range contacts {
		go rt.update(contact)
	}
}

func CreateTestRT() ([]Contact, *RoutingTable) {
	rt := newRoutingTable(NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	contacts = append(contacts, NewContact(
		NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))

	for _, contact := range contacts {
		rt.update(contact)
	}
	return contacts, rt
}

func CreateTestRT2() ([]Contact, *RoutingTable) {
	rt := newRoutingTable(NewContact(
		NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	contacts = append(contacts, NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8019"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	contacts = append(contacts, NewContact(
		NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))

	for _, contact := range contacts {
		rt.update(contact)
	}
	return contacts, rt
}

func CreateTestRT3() ([]Contact, *RoutingTable) {
	rt := newRoutingTable(NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8007"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8006"))
	contacts = append(contacts, NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8007"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	contacts = append(contacts, NewContact(
		NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))

	for _, contact := range contacts {
		rt.update(contact)
	}
	return contacts, rt
}

func CreateTestRT4() ([]Contact, *RoutingTable) {
	rt := newRoutingTable(NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8009"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8008"))
	contacts = append(contacts, NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8009"))

	for _, contact := range contacts {
		rt.update(contact)
	}
	return contacts, rt
}

func CreateTestRT10() ([]Contact, *RoutingTable) {
	rt := newRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:9500"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:9501"))

	for _, contact := range contacts {
		rt.update(contact)
	}
	return contacts, rt
}

func CreateTestRT11() ([]Contact, *RoutingTable) {
	rt := newRoutingTable(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:9501"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:9500"))

	for _, contact := range contacts {
		rt.update(contact)
	}
	return contacts, rt
}
