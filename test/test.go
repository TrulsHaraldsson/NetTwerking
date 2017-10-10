package main

import (
	"fmt"
)

func main() {
	/*reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	*/
	text := "hej"
	byteArray := []byte(text)
	stringFromByte := string(byteArray)
	byteAgain := []byte(stringFromByte)

	fmt.Println("string:", text, "byte:", byteArray, "String from byte:", stringFromByte)
	fmt.Println("byte2:", byteAgain)
}
