package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	clientAddrSTruct, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal("error while making udp socket")
	}
	//label the data packet when sending with my server socket and the client socket
	conn, err := net.DialUDP("udp", nil, clientAddrSTruct)
	defer conn.Close()

	if err != nil {
		log.Fatal("error while esstablih connection")
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("write what you want to send over our upd socket!!!")
	for {
		fmt.Print(">")
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("erro while reading the msg ,retype it ....")
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("data didn't reached / sending eror")
		}
		fmt.Println("mesage ssent!!")
	}

}
