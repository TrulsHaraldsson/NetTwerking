package d7024e

import (
	"crypto/sha256"
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

func NewValueID(filename *string) *KademliaID {
	hash := sha256.Sum256([]byte(*filename))
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = hash[i]
	}
	return &newKademliaID
}

/*func NewKademliaIDFromByteArray(id []byte) *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = id[i]
	}
	return &newKademliaID
}*/

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

func NewRandomKademliaIDWithPrefix(prefix KademliaID, index int) *KademliaID {
	rand := NewRandomKademliaID()
	prefixOne := &KademliaID{}
	bitsFlipped := 0
	for i := 0; i < IDLength; i++ {
		byte := prefixOne[i]
		var i2 uint
		for i2 = 0; i2 < 8; i2++ {
			if bitsFlipped < index {
				byte = byte | (1 << (7 - i2))
				bitsFlipped += 1
			} else {
				break
			}
		}
		prefixOne[i] = byte
	}
	prefixZero := NewKademliaID("ffffffffffffffffffffffffffffffffffffffff")
	bitsFlipped = 0
	for i := 0; i < IDLength; i++ {
		byte := prefixZero[i]
		for i2 := 0; i2 < 8; i2++ {
			if bitsFlipped < index {
				byte = byte >> 1
				bitsFlipped += 1
			} else {
				break
			}
		}
		prefixZero[i] = byte
	}
	for i3 := 0; i3 < IDLength; i3++ {
		rand[i3] = rand[i3] & prefixZero[i3]
		prefix[i3] = prefix[i3] & prefixOne[i3]
		rand[i3] = rand[i3] | prefix[i3]
	}
	return rand
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
