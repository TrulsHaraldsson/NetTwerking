package d7024e

import (
	"testing"
)

func TestKademliaIDUniqness(t *testing.T) {
	var kID1 = NewRandomKademliaID()
	var kID2 = NewRandomKademliaID()
	if kID1 == kID2 {
		t.Error("Expected two unique kademliaIDs, got ID1 = ", kID1,
			" and ", kID2)
	}
}

func TestKademliaIDLength(t *testing.T) {
	var kademliaID = NewRandomKademliaID()
	if len(*kademliaID) != IDLength {
		t.Error("Expected length", IDLength, ", got ", len(*kademliaID))
	}
}

func TestKademliaIDByString(t *testing.T) {
	var sID1 = "ffffffffffffffffffffffffffffffffffffffff"
	var kID1 = NewKademliaID(sID1)

	if kID1.String() != "ffffffffffffffffffffffffffffffffffffffff" {
		t.Error(kID1.String())
	}
}

func TestKademliaIDLess(t *testing.T) {
	var sID1 = "ffffffffffffffffffffffffffffffffffffffff"
	var sID2 = "0000000000000000000000000000000000000000"
	var kID1 = NewKademliaID(sID1)
	var kID2 = NewKademliaID(sID2)
	if kID1.Less(kID2) {
		t.Error("Expected", kID1, " to be less than ", kID2)
	}

	var sID3 = "ffffffffffffffffffffffffffffffffffffffff"
	var kID3 = NewKademliaID(sID3)
	if kID1.Less(kID3) {
		t.Error("Expected", kID1, " not to be less than ", kID3)
	}
}

func TestKademliaIDEquals(t *testing.T) {
	var sID1 = "0000000000000000000000000000000000000000"
	var sID2 = "0000000000000000000000000000000000000000"
	var kID1 = NewKademliaID(sID1)
	var kID2 = NewKademliaID(sID2)
	if !kID1.Equals(kID2) {
		t.Error("Expected following IDs to be equal \nID1 = ", sID1, "\nID2 = ", sID2)
	}
}

func TestKademliaIDCalcDistance(t *testing.T) {
	var sID1 = "0000000000000000000000000000000000000010"
	var sID2 = "0000000000000000000000000000000000000001"
	var expected = "0000000000000000000000000000000000000011"
	var kID1 = NewKademliaID(sID1)
	var kID2 = NewKademliaID(sID2)
	var distance = kID1.CalcDistance(kID2)

	if *distance != *NewKademliaID(expected) {
		t.Error("Expected distance to be", expected, ", got", distance)
	}
}
