package main

import (
	"fmt"

	"../golang"
)

func main() {
	/*reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	*/
	s1 := "lol"
	s2 := "lol2"
	a := d7024e.NewValueID(&s1).String()
	b := d7024e.NewValueID(&s2).String()
	fmt.Println(a, "\n"+b)
}
