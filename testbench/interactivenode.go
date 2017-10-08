package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"strconv"
	"time"
	"reflect"
	"io/ioutil"
	"net"
	"../golang"
)

/*
 * INSTRUCTIONS TO START:
 * 1. sudo docker build -t server
 * 2. sudo docker run server
 * NOTE: Now the server will run indefinitely.
 *
 * INSTRUCTIONS TO STOP:
 * 1. Open a second terminal
 * 2. sudo docker container ls
 * 3. find 'container id' for the server
 * 3. sudo stop 'container id'
 */
const MESSAGE_SIZE = 1024
const PORT = 7999

func main() {
	// (flagName, default value, description)
	//address := flag.String("addr", "127.0.0.1", "IP-address of this node")
	port := flag.Int("port", 65000, "Port this node will listen to")

	interactive := flag.String("interactive", "false", "Manual Control")

	id := flag.String("id", "none", "id of node to connect to")
	addr := flag.String("addr", "none", "ip address of node to connect to")

	flag.Parse()

	fmt.Println(*port)

	if *interactive == "true" {
		fmt.Println("Booting up an interactive node")
	} else {
		fmt.Println("Booting up a non-interactive node")
	}

	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}

	fmt.Println("Host name:", name)
	addrs, err := net.LookupHost(name)

	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}
	for _, a := range addrs {
		fmt.Println("ADDRESS", a)

	}

	time.Sleep(1 * time.Second)

	reader := bufio.NewReader(os.Stdin)
	if *interactive == "true" {
<<<<<<< HEAD
		kID := "abcdef1234abcdef1234abcdef1234abcdef1234"
		kademlia := d7024e.CreateAndStartNode(addrs[0]+":"+strconv.Itoa(*port), kID, nil)
		//kademlia := d7024e.CreateAndStartNode(addrs[0]+":"+strconv.Itoa(*port), "none", "none")
		fmt.Println("You can always type TIPS, to see different commands!")
		tips()
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]
		for {
			if input == "QUIT" {
				break
			}
			if input == "TIPS"{
				tips()
			}
			switch input {
			case d7024e.PING: //TODO: Dont assume localhost
				onPing2(kademlia, reader)
			case d7024e.FIND_NODE:
				onFindNode(kademlia, reader)
			case d7024e.FIND_VALUE:
				onFindValue(kademlia, reader)
			case d7024e.STORE:
				fmt.Println("Need help to fix with /app but everything works, just cant see stored files in the folder files.")
				onStore(kademlia, reader)
			case "ls": //shows 20 closest contacts to a random id.
				c := d7024e.NewContact(d7024e.NewKademliaID(kID), "lol")
				fmt.Println(kademlia.LookupContact(&c))
			default:
				fmt.Println("Wrong syntax in message, ignoring it...")
			}
			input, _ = reader.ReadString('\n')
			input = input[:len(input)-1]
		}
	} else {
		time.Sleep(1 * time.Second)
		fmt.Println("Connecting to addr:", *addr)
		fmt.Println("With ID:", *id)
		initContact := d7024e.NewContact(d7024e.NewKademliaID(*id), *addr)
		d7024e.CreateAndStartNode(addrs[0]+":"+strconv.Itoa(*port), "none", &initContact)
		for {
			time.Sleep(1 * time.Second)
		}
		//fmt.Println("Press enter to quit")
		//_, _ = reader.ReadString('\n')
	}

}
func tips(){
	fmt.Println("\nFollowing options are valid:")
	fmt.Println("QUIT, quits the program.")
	fmt.Println("PING, send ping message to a given port on localhost.")
	fmt.Println("FIND_NODE, search for closest nodes to specified ID.")
	fmt.Println("STORE, store a file with a given name")
	fmt.Println("FIND_VALUE, find a file by its name")
//		fmt.Println("LES, see Directory of files")
}
func onPing(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	fmt.Println("Please write a port to send to e.g. 65002.")
	rport, _ := reader.ReadString('\n')
	rport = rport[:len(rport)-1]
	rp, _ := strconv.Atoi(rport)
	ok := kademlia.Ping(d7024e.CreateAddr("127.0.0.1", rp))
	fmt.Println("ping status:", ok)
}

func onPing2(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	fmt.Println("Please write a address to send to e.g. 'localhost:65002'.")
	raddr, _ := reader.ReadString('\n')
	raddr = raddr[:len(raddr)-1]
	ok := kademlia.Ping(raddr)
	fmt.Println("ping status:", ok)
}

func onFindNode(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	fmt.Println("Please write a kademliaID to search for. 'none', for a random id")
	rid, _ := reader.ReadString('\n')
	rid = rid[:len(rid)-1]
	var kID *d7024e.KademliaID
	if rid == "none" {
		kID = d7024e.NewRandomKademliaID()
	} else {
		kID = d7024e.NewKademliaID(rid)
	}
	contacts := kademlia.SendFindContactMessage(kID)
	fmt.Println("Number of contacts found:", len(contacts))
}

/* Old function
func onStore(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	fmt.Println("Please write the path/name of the file you want to store.")
	rid, _ := reader.ReadString('\n')
	f, _ := os.Open(rid[:len(rid)-1])
	content, _ := ioutil.ReadAll(f)
	content = []byte(content)
	fmt.Println("You wish to store \n", string(content), "\nUnder which name?")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1]
	kademlia.SendStoreMessage(&name, &content)
	fmt.Println("Store Sent!", name)
}
*/

func new(kademlia *d7024e.Kademlia, reader *bufio.Reader){
	fmt.Println("Write the name of the new file")
	filename, _ := reader.ReadString('\n')
	filename = filename[:len(filename)-1]
	fmt.Println("Write content of file")
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	data := []byte(text)

	file := "./files/" + filename

	err := ioutil.WriteFile(file, data, 0644)
	if err != nil{
		panic(err)
	}

	content, err2 := ioutil.ReadFile(file)
	if err2 != nil {
		fmt.Println("dont exist: ", err2)
	}

	valueID := kademlia.SendStoreMessage(&filename, &content)
	//valueID := kademlia.SendStoreMessage(&filename, &data)
	if valueID != nil{
		//fmt.Println("Store successful and returned filename (type) : ", reflect.TypeOf(valueID),"\n",  valueID)
		fmt.Println("Storemessage successful!")
		//-- direct search on the stored file --//
/*		fmt.Println("\nWrite the name again to see content.\n")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		fmt.Println("Name : ", text, reflect.TypeOf(text))
		find := kademlia.SendFindValueMessage(&text)
		if find == nil {
			fmt.Println("Not found!")
		}else{
			file := string(find)
			fmt.Println("Returned file content : ", string(file), "type : ", reflect.TypeOf(file))
		}*/
	}else{ //-- The storeMessage failed --//
		fmt.Println("Storemessage unsuccessful!")
	}
}

func old(kademlia *d7024e.Kademlia, reader *bufio.Reader){
	fmt.Println("Write the name of the old loaded file or see files by typing 'ls'")
	filename, _ := reader.ReadString('\n')
	filename = filename[:len(filename)-1]
	file := string(filename)
	if file != "ls" {

		//-- Read the old file --//

		name := "./files/" + file
		content, err2 := ioutil.ReadFile(name)
		if err2 != nil {
			fmt.Println("dont exist: ", err2)
		}
		fmt.Println("File contents: ", string(content),"\n")

		//-- send the file and get it --//

		valueID := kademlia.SendStoreMessage(&file, &content)
		if valueID != nil{
			//fmt.Println("Store successful and returned filename (type) : ", reflect.TypeOf(valueID),"\n",  valueID)
			fmt.Println("Storemessage successful!")
	/*
			fmt.Println("\nWrite the name again to see content.\n")
			text, _ := reader.ReadString('\n')
			text = text[:len(text)-1]
			fmt.Println("Name : ", text, reflect.TypeOf(text))
			find := kademlia.SendFindValueMessage(&text)
			if find == nil {
				fmt.Println("Not found!")
			}else{
				file := string(find)
				fmt.Println("Returned file content : ", string(file), "type : ", reflect.TypeOf(file))
			}*/
		}else{
			fmt.Println("Storemessage unsuccessful!")
		}
	}else{
		onDirectory(kademlia, reader)
	}
}

func onStore(kademlia *d7024e.Kademlia, reader *bufio.Reader){
	fmt.Println("Want to store a new file or load an already existing one?\n new or old")
	choice, _ := reader.ReadString('\n')
	choice = choice[:len(choice)-1]
	if choice == "new" {
		//-- User creates new file that will be sent out to network --//
		new(kademlia, reader)
	}else{
		//-- When a user wants to load in an old file --//
		old(kademlia, reader)
	}
	fmt.Println("Done\n")
}

func onFindValue(kademlia *d7024e.Kademlia, reader *bufio.Reader){
	fmt.Println("\nWrite the name of the file you wish to see content off.\n")
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	//fmt.Println("Name : ", text, reflect.TypeOf(text))
	find := kademlia.SendFindValueMessage(&text)
	if find == nil {
		fmt.Println("Not found!")
	}else{
		file := string(find)
		fmt.Println("Returned File content : ", string(file))//, "type : ", reflect.TypeOf(file))
	}
	fmt.Println("Done\n")
}

func onDirectory(kademlia *d7024e.Kademlia, reader *bufio.Reader){
	path := "./files"
	files, err := ioutil.ReadDir(path)
  if err != nil {
      log.Fatal(err)
  }
	fmt.Println("\n")
  for _, f := range files {
          fmt.Println(f.Name())
	}
}
