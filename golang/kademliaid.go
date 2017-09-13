package d7024e

import (
	"encoding/hex"
	"math/rand"
	"time"
)


const IDLength = 20

type KademliaID [IDLength]byte

/*
 * Reads a string and returns a KademliaID. The function assumes that the string
 * resembles a hexdecimal number.
 */
func NewKademliaID(data string) *KademliaID {
	decoded, _ := hex.DecodeString(data)

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}

	return &newKademliaID
}

/*
 * Creates a new random-object with a seed based on current time. Hopefully this will be 
 * enough to create different seeds for all threads.
 */
func NewRandomKademliaID() *KademliaID {
	var t = time.Now().UnixNano()
	var r = rand.New(rand.NewSource(t))
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(r.Intn(256))
	}
	return &newKademliaID
}

/*
 * Check if the calling KademliaID is less than the given KademliaID 
 */ 
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

/*
 * Returns true if the two KademliaIDs are identical
 */
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

/*
 * Calculates the XOR distance between two KademliaIDs. The result will be a new
 * KademliaID.
 */
func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

/*
 * Returns a string representation of the KademliaID
 */
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}
