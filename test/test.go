package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	/*reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	*/
	content, err := ioutil.ReadFile("dahdkefsfke")
	fmt.Println("content:", content, "err:", err)
}
