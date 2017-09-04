package d7024e

import (
	"fmt"
	"testing"
)


func TestNewRandomKademliaID(t *testing.T) {
	var kademliaid = NewRandomKademliaID()
	fmt.Println(kademliaid)
}
