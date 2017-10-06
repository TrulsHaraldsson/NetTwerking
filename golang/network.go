package d7024e

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"
)

const MESSAGE_SIZE = 1024

type Network struct {
	alpha    int
	kademlia *Kademlia
	addr     string
}

func NewNetwork(alpha int, addr string) Network {
	return Network{alpha: alpha, addr: addr}
}

/*
* Starts a UDP socket listening on port and ip specified.
* When a package is received it will start a new thread handling it.
 */
func (network Network) Listen() {
	//addrServer := CreateAddr(ip, port)
	udpConn, err := net.ListenPacket("udp", network.addr)
	if err != nil {
		panic(err)
	} else {
		//fmt.Println("Listening to port", port)
	}
	defer udpConn.Close()
	for {
		b, addrClient, err := network.readAnswer(udpConn)
		//fmt.Println("message received.")
		if err != nil {
			// handle error
			//fmt.Println("Error when reading from socket...", err)
		} else {
			m, mData, err2 := UnmarshallMessage(b)
			if err2 != nil {
				fmt.Println("Error when unmarshalling message...", err2)
			} else {
				go network.handleConnection(m, mData, addrClient)
				//fmt.Println("Starting new thread to handle connection...")
			}
		}
	}
}

/* Function to handle an incoming package. Assuming package is a serialized struct of type "Message".
* Depending on the message type, different functions are run.
* Also gives a answer back to the caller.
 */
func (network Network) handleConnection(message Message, mData interface{}, addr net.Addr) {
	switch message.MsgType {
	case PING:
		//fmt.Println("Ping message received")
		//fmt.Println("address is:", addr.String())
		network.kademlia.OnPingMessageReceived(&message, addr)
	case FIND_NODE:
		//fmt.Println("Find node message received")
		network.kademlia.OnFindNodeMessageReceived(&message, mData.(FindNodeMessage), addr)
	case FIND_VALUE:
		//fmt.Println("Find value message received")
		network.kademlia.OnFindValueMessageReceived(&message, mData.(FindValueMessage), addr)
	case STORE:
		//fmt.Println("Store message received")
		network.kademlia.OnStoreMessageReceived(&message, mData.(StoreMessage), addr)

	default:
		fmt.Println("Wrong syntax in message, ignoring it...")
		return
	}
	network.kademlia.RT.update(message.Sender)
}

/*
* Creates a string address
 */
func CreateAddr(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

/*
 * Sends a ping to given address
 */
func (network *Network) SendPingMessage(addr string, msg *Message) (Message, error) {
	response, _, err := network.sendSpecificMessage(addr, msg, PING_ACK)
	return response, err
}

/*
 * Send a find contact message to the given addr.
 * Returns the response of from the given addr.
 */
func (network *Network) sendFindContactMessage(addr string, msg *Message) (Message, AckFindNodeMessage, error) {
	response, responseData, err := network.sendSpecificMessage(addr, msg, FIND_NODE_ACK)
	if err != nil {
		return response, AckFindNodeMessage{}, err
	}
	return response, responseData.(AckFindNodeMessage), err
}

/*
 * Request to find a value over the network.
 */
func (network *Network) sendFindValueMessage(addr string, msg *Message) (Message, AckFindValueMessage, error) {
	response, responseData, err := network.sendSpecificMessage(addr, msg, FIND_VALUE_ACK)
	if err != nil {
		return response, AckFindValueMessage{}, err
	}
	return response, responseData.(AckFindValueMessage), err
}

/*
 * Sends a message over the network to the alpha closest neighbors in the routing table and waits for response
 * from neighbor OnStoreMessageReceived func.
 */
func (network *Network) sendStoreMessage(addr string, msg *Message) (Message, error) {
	response, _, err := network.sendSpecificMessage(addr, msg, STORE_ACK)
	return response, err
}

/*
 * Sends a message of specific type defind by Message.
 * Returns the Ack type of the given Message. e.g. Message = PING, Ack = PING_ACK 
 */
func (network *Network) sendSpecificMessage(addr string, msg *Message, responseType string) (Message, interface{}, error) {
	response, rData, err := network.sendMessage(addr, *msg)
	if err != nil {
		return response, rData, err
	}
	if msg.RPC_ID == response.RPC_ID {
		if response.MsgType == responseType {
			return response, rData, nil
		} else {
			return Message{}, rData, errors.New("Wrong message sent back, it is not a FindNodeAck...")
		}
	} else {
		return Message{}, rData, errors.New("Wrong RPC_ID returned, it is not from the server, sent to...")
	}
}

/*
 * Sends a message of type "Message" to address specified.
 * Waits for a response and unmarshalls it as a Message and the MessageData type.
 * Returns both Message and MessageData
 */
func (network *Network) sendMessage(addr string, message Message) (Message, interface{}, error) {
	msgJson, err := MarshallMessage(message)
	if err != nil {
		return Message{}, nil, err
	}
	b, err2 := network.sendData(addr, msgJson)
	if err2 != nil {
		return Message{}, nil, err2
	}
	returnMsg, msgData, err3 := UnmarshallMessage(b)
	return returnMsg, msgData, err3
}

/*
 * Sends data to the given address. Then waits for a response or timeout.
 * Returns: The response as a byte array.
 */
func (network *Network) sendData(addr string, data []byte) ([]byte, error) {
	var returnMsg []byte
	
	ip := regexp.MustCompile(":").Split(network.addr, 2)[0] //Take port and convert to int
	addrLocal := CreateAddr(ip, 0)
	addrRemote, _ := net.ResolveUDPAddr("udp", addr)
	udpConn, err := net.ListenPacket("udp", addrLocal)
	if err != nil {
		fmt.Println("network/SendData: err ", err)
		return returnMsg, err
	}

	// Waits for timeOut until execption is thrown.
	timeOut := 1 * time.Second 
	udpConn.SetDeadline(time.Now().Add(timeOut))
	defer udpConn.Close()
	_, err2 := udpConn.WriteTo(data, addrRemote)
	if err2 != nil {
		fmt.Println("network/SendData: err2 ", err2)
		return returnMsg, err2
	}

	returnMsg, _, err3 := network.readAnswer(udpConn)
	if err3 != nil {
		fmt.Println("network/SendData: err3 ", err3)
		return returnMsg, err3
	}

	return returnMsg, nil
}

/*
 * Reads from the specified UDP net.PacketConn 
 * Returns the data as a byte array and the adress from where it came.
 */
func (network *Network) readAnswer(udpConn net.PacketConn) ([]byte, net.Addr, error) {
	b := make([]byte, MESSAGE_SIZE)
	n, addr, err := udpConn.ReadFrom(b)
	b = b[:n]
	if err != nil {
		return b, addr, err
	}
	return b, addr, nil
}

/*
 * Writes a message of type Message to addr, does not wait for response.
 */
func (network *Network) WriteMessage(addr string, message Message) error {
	msgJson, err := MarshallMessage(message)
	if err != nil {
		return err
	}
	err2 := network.connectAndWrite(addr, msgJson)
	return err2
}

/*
 * Connects to addr and writes the message.
 * Does not wait for a response.
 */
func (network *Network) connectAndWrite(addr string, message []byte) error {
	//addrLocal := CreateAddr("localhost", 0)
	ip := regexp.MustCompile(":").Split(network.addr, 2)[0] //Take port and convert to int
	addrLocal := CreateAddr(ip, 0)
	addrRemote, _ := net.ResolveUDPAddr("udp", addr)
	udpConn, err := net.ListenPacket("udp", addrLocal)
	if err != nil {
		return err
	}
	defer udpConn.Close()
	_, err2 := udpConn.WriteTo(message, addrRemote)
	if err2 != nil {
		return err2
	} else {
		return nil
	}
}
