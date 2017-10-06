package main

import (
	//"bufio"
	"flag"
	"fmt"
	"os"
	//"strconv"
	//"time"
	"net"
	//"../golang"
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
	
	interactive := flag.Bool("interactive", false, "Manual Control")

	flag.Parse()

	fmt.Println(*port)

	if *interactive == true {
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

/*	for {
		fmt.Println("Running", time.Now())
		time.Sleep(1 * time.Second)
	}
	//fmt.Println("IP:", *address)
	//fmt.Println("PORT:", *port)

	kID := d7024e.NewRandomKademliaID()
	c := d7024e.NewContact(kID, "localhost:port")
	fmt.Println("Node ID:", kID)

	kademlia := d7024e.NewKademlia(*port, "none")
	//contact := d7024e.NewContact(kID, d7024e.CreateAddr(*address, *port))
	//rt := d7024e.NewRoutingTable(contact)
	//kademlia := d7024e.Kademlia{rt, 20}
	//network := d7024e.NewNetwork(3, kademlia)
	go kademlia.Start()

	time.Sleep(1 * time.Second)

	reader := bufio.NewReader(os.Stdin)
	if *interactive == true {
		fmt.Println("Following options are valid:")
		fmt.Println("QUIT, quits the program.")
		fmt.Println("PING, send ping message to a given port on localhost.")
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]
		for {
			if input == "QUIT" {
				break
			}
			switch input {
			case d7024e.PING:
				fmt.Println("Please write a port to send to e.g. 65002.")
				rport, _ := reader.ReadString('\n')
				rport = rport[:len(rport)-1]
				rp, _ := strconv.Atoi(rport)
				msg := d7024e.NewPingMessage(&c)
				response, _, _ := d7024e.SendMessage(
					d7024e.CreateAddr("127.0.0.1", rp), msg)
				fmt.Println(response)
			default:
				fmt.Println("Wrong syntax in message, ignoring it...")
			}
			input, _ = reader.ReadString('\n')
			input = input[:len(input)-1]
		}
	} else {
		fmt.Println("Press enter to quit")
		_, _ = reader.ReadString('\n')
	}
*/
}
