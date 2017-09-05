package d7024e

import (
	"testing"
	"fmt"
	"encoding/json"
)

func TestLookupContact(t *testing.T) {

}

func spawnCh() chan string{	
	fmt.Println("Creating channels")
	ch := make(chan string,2)
	go func(){
		ch <- "hello"
	}()
	return ch
}

func TestStoreItems(t *testing.T){
	fmt.Println("Testing multiple stores")

	ch1 := spawnCh()
	ch2 := spawnCh()
	
	for i := 0; i < 2; i++{
		select {
			case n := <- ch1: 
				fmt.Printf("ch1 : %s\n", n)
				data := []byte(n)
				
				var kademlia Kademlia
				kademlia.Store(data)
				kademlia.GetList()		
				
				
			case n := <- ch2:
				fmt.Printf("ch2 : %s\n", n)
				data := []byte(n)
				
				var kademlia Kademlia
				kademlia.Store(data)
				kademlia.GetList()
		}
	}
}