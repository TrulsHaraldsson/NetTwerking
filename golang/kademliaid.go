package d7024e

import (
	"encoding/hex"
	"math/rand"
	"time"
)


const IDLength = 20

type KademliaID [IDLength]byte

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
	var r = rand.New(rand.NewSource(time.Now().Unix()))
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(r.Intn(256))
	}
	return &newKademliaID
}

func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}
