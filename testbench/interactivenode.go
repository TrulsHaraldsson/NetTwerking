package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
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
		kID := "abcdef1234abcdef1234abcdef1234abcdef1234"
		kademlia := d7024e.CreateAndStartNode(addrs[0]+":"+strconv.Itoa(*port), kID, nil)
		//kademlia := d7024e.CreateAndStartNode(addrs[0]+":"+strconv.Itoa(*port), "none", "none")
		fmt.Println("You can always type HELP, to see different commands!")
		help()
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]
		for {
			if input == "QUIT" {
				break
			}
			if input == "HELP" || input == "help" {
				help()
			}
			switch input {
			case d7024e.PING: //TODO: Dont assume localhost
				onPing(kademlia, reader)
			case d7024e.FIND_NODE:
				onFindNode(kademlia, reader)
			case d7024e.FIND_VALUE:
				onFindValue(kademlia, reader)
			case d7024e.STORE:
				onStore(kademlia, reader)
			case "DELETE":
				onDelete(kademlia, reader)
			case "dir":
				onDirectory(kademlia, reader)
			case "rt": //shows 20 closest contacts to a random id.
				fmt.Println(kademlia.RT.Contacts())
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
func help() {
	fmt.Println("\nFollowing options are valid:")
	fmt.Println("QUIT, quits the program.")
	fmt.Println("PING, send ping message to a given port on localhost.")
	fmt.Println("FIND_NODE, search for closest nodes to specified ID.")
	fmt.Println("STORE, store a file with a given name")
	fmt.Println("FIND_VALUE, find a file by its name")
	fmt.Println("DELETE, delete a file form node")
	fmt.Println("dir, list all files")
	fmt.Println("rt, show how many contacts in routing table")
}

func onPing(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	fmt.Println("Please write a address to ping, below are some choices.")
	newContact := d7024e.NewContact(d7024e.NewRandomKademliaID(), "no address")
	contacts := kademlia.LookupContact(&newContact)
	for i := 0 ; i < 10 ; i ++ {
		fmt.Println("Contact : ", contacts[i].Address)
	}

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

func new(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	fmt.Println("Write the name of the new file")
	filename, _ := reader.ReadString('\n')
	filename = filename[:len(filename)-1]
	fmt.Println("Write content of file")
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	data := []byte(text)

	file := "../newfiles/" + filename

	err := ioutil.WriteFile(file, data, 0644)
	if err != nil {
		panic(err)
	}
	content, err2 := ioutil.ReadFile(file)
	if err2 != nil {
		fmt.Println("dont exist: ", err2)
	}

	valueID := kademlia.SendStoreMessage(&filename, &content)
	if valueID == nil {
		fmt.Println("Storemessage successful!")
	} else { //-- The storeMessage failed --//
		fmt.Println("Storemessage unsuccessful!")
	}
}

func old(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	fmt.Println("Write the name of the old loaded file or see files by typing 'dir'")
	filename, _ := reader.ReadString('\n')
	filename = filename[:len(filename)-1]
	file := string(filename)
	if file != "dir" {

		//-- Read the old file --//

		name := "../newfiles/" + file
		content, err2 := ioutil.ReadFile(name)
		if err2 != nil {
			fmt.Println("dont exist: ", err2)
		}else{
			fmt.Println("File contents: ", string(content), "\n")

			//-- send the file and get it --//

			valueID := kademlia.SendStoreMessage(&file, &content)
			if valueID == nil {
				fmt.Println("Storemessage successful!")
			}else{
					fmt.Println("Storemessage unsuccessful!")
			}
		}
	}else{
		onDirectory(kademlia, reader)
	}
}

func onStore(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
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
}

func onFindValue(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	fmt.Println("\nWrite the name of the file you wish to see content off.\n")
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	find := kademlia.SendFindValueMessage(&text)
	if find == nil {
		fmt.Println("Not found!")
	}else{
		file := string(find)
		fmt.Println("Returned File content : ", string(file))
		err := ioutil.WriteFile("../newfiles/"+text, find, 0644)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Done\n")
}

func onDirectory(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	path := "../newfiles"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n")
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func onDelete(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	fmt.Println("\nWrite the name of the file you wish to delete.\n")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1]
	kademlia.DeleteFile("../newfiles/" + name)
}
