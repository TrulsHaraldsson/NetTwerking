package d7024e

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	_, rt := CreateTestRT()

	contacts := rt.FindClosestContacts(
		NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}


func TestRoutingTableAddContact(t *testing.T) {
	rt := NewRoutingTable(NewContact(
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
		go rt.AddContact(contact)
	}
}

func CreateTestRT() ([]Contact, *RoutingTable) {
	rt := NewRoutingTable(NewContact(
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
		rt.AddContact(contact)
	}
	return contacts, rt
}

func CreateTestRT2() ([]Contact, *RoutingTable) {
	rt := NewRoutingTable(NewContact(
		NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	contacts = append(contacts, NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	contacts = append(contacts, NewContact(
		NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))

	for _, contact := range contacts {
		rt.AddContact(contact)
	}
	return contacts, rt
}


func CreateTestRT3() ([]Contact, *RoutingTable) {
	rt := NewRoutingTable(NewContact(
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
		rt.AddContact(contact)
	}
	return contacts, rt
}

func CreateTestRT4() ([]Contact, *RoutingTable) {
	rt := NewRoutingTable(NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8009"))
	contacts := []Contact{}
	contacts = append(contacts, NewContact(
		NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8008"))
	contacts = append(contacts, NewContact(
		NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8009"))

	for _, contact := range contacts {
		rt.AddContact(contact)
	}
	return contacts, rt
}
