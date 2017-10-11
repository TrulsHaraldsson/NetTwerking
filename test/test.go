package main

import (
	"time"

	"../golang"
)

func main() {
	/*reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	*/
	d7024e.CreateAndStartNode("localhost:8003", "none", nil)

	k2 := d7024e.CreateAndStartNode("localhost:8004", "none", nil)

	time.Sleep(40 * time.Millisecond)
	k2.Ping("localhost:8003")
	time.Sleep(40 * time.Millisecond)
	for i := 0; i < 1000; i++ {
		k2.FindContact(d7024e.NewRandomKademliaID())
	}
}
