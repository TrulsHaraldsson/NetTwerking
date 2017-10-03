package main

import (
	"fmt"

	"../golang"
)

func main() {
	id := d7024e.NewKademliaID("4f00000000000000000000000000000000000000")
	rand := d7024e.NewRandomKademliaIDWithPrefix(*id, 5)
	fmt.Println(rand)
	fmt.Println(id)
}
