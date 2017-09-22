package main

import (
	"fmt"
	"time"
)

var variable int = 5

func main() {
	ch := make(chan int)  // index, newValue
	ch2 := make(chan int) // for returning
	go callAlterSomething(ch)
	go callAlterSomething(ch)
	go callGetSomething(ch, ch2)
	go callGetSomething(ch, ch2)
	go callGetSomething(ch, ch2)
	go callAlterSomething(ch)
	go callAlterSomething(ch)
	go callGetSomething(ch, ch2)
	go callAlterSomething(ch)
	go callGetSomething(ch, ch2)
	go callGetSomething(ch, ch2)
	go callGetSomething(ch, ch2)
	for {
		x := <-ch
		switch x {
		case 1:
			alterSomething()
		case 2:
			ch2 <- getSomething()
		}
	}

}

func callAlterSomething(ch chan int) {
	fmt.Println("calling alter")
	ch <- 1
	fmt.Println("alter call finished")
}

func callGetSomething(chOut, chIn chan int) {
	fmt.Println("calling get")
	chOut <- 2
	x := <-chIn
	fmt.Println("received:", x, "get call finished")
}

func getSomething() int {
	fmt.Println("starting getting")
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("ending getting")
	return variable
}

func alterSomething() {
	fmt.Println("starting alter")
	time.Sleep(1000 * time.Millisecond)
	variable += 1
	fmt.Println("ending alter")
}
