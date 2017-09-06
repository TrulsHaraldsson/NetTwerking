package main

import (
	"fmt"
	"flag"
	"bufio"
	"strconv"
	"os"
	"time"
	"../golang"
)

const MESSAGE_SIZE = 1024

func main() {
	// (flagName, default value, description)
	address := flag.String("addr", "127.0.0.1", "IP-address of this node")
	port := flag.Int("port", 65000, "Port this node will listen to")
	interactive := flag.Bool("interactive", false, "Manual Control")
	
	flag.Parse()

	if *interactive == true {
		fmt.Println("Booting up an interactive node")
	}else {
		fmt.Println("Booting up a non-interactive node")
	}
	
	fmt.Println("IP:",*address)
	fmt.Println("PORT:",*port)
	
	kID := d7024e.NewRandomKademliaID()
	fmt.Println("Node ID:", kID)

	contact := d7024e.NewContact(kID, d7024e.CreateAddr(*address, *port))
	rt := d7024e.NewRoutingTable(contact)
	kademlia := d7024e.Kademlia{rt, 20}
	network := d7024e.NewNetwork(3, kademlia)
	go network.Listen(*address, *port)
	
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
				break;
			}
			switch input {
			case d7024e.PING:
				fmt.Println("Please write a port to send to e.g. 65002.")
				rport, _ := reader.ReadString('\n')
				rport = rport[:len(rport)-1]
				rp, _ := strconv.Atoi(rport)
				msg := d7024e.NewPingMessage(kID)
				response, _ := d7024e.SendMessage(
					d7024e.CreateAddr("127.0.0.1", rp), msg)
				fmt.Println(response)
			default:
				fmt.Println("Wrong syntax in message, ignoring it...")
			}
			input, _ = reader.ReadString('\n')
			input = input[:len(input)-1]
		}
	}else {
		fmt.Println("Press enter to quit")
		_, _ = reader.ReadString('\n')
	}

}

