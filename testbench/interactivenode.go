package main

import (
	//"bufio"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
	//"strconv"

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
		kademlia := d7024e.CreateAndStartNode(addrs[0]+":"+strconv.Itoa(*port), "none", "none")
		fmt.Println("Following options are valid:")
		fmt.Println("QUIT, quits the program.")
		fmt.Println("PING, send ping message to a given port on localhost.")
		fmt.Println("FIND_NODE, search for closest nodes to specified ID.")
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]
		for {
			if input == "QUIT" {
				break
			}
			switch input {
			case d7024e.PING: //TODO: Dont assume localhost
				onPing2(kademlia, reader)
			case d7024e.FIND_NODE:
				onFindNode(kademlia, reader)
			case d7024e.FIND_VALUE:
				onFindValue(kademlia, reader)
			case d7024e.STORE:
				onStore(kademlia, reader)
			default:
				fmt.Println("Wrong syntax in message, ignoring it...")
			}
			input, _ = reader.ReadString('\n')
			input = input[:len(input)-1]
		}
	} else {
		time.Sleep(1 * time.Second)
		d7024e.CreateAndStartNode(addrs[0]+":"+strconv.Itoa(*port), "none", "10.0.0.2:7999")
		for {
			time.Sleep(1 * time.Second)
		}
		//fmt.Println("Press enter to quit")
		//_, _ = reader.ReadString('\n')
	}

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
	fmt.Println("Contacts found:", contacts)
}

func onFindValue(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	fmt.Println("Please write the kademliaID that the VALUE have")
	/*rid, _ := reader.ReadString('\n')
	rid = rid[:len(rid)-1]
	var kID *d7024e.KademliaID
	if rid == "none" {
		kID = d7024e.NewRandomKademliaID()
	} else {
		kID = d7024e.NewKademliaID(rid)
	}
	data_as_byte_array := kademlia.SendFindValueMessage(kID)
	fmt.Println("Value found:", data_as_byte_array)
	*/

}

func onStore(kademlia *d7024e.Kademlia, reader *bufio.Reader) {
	fmt.Println("Not implemented yet")
}
