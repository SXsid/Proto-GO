package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SXsid/ProtoGo/internal/request"
)

func main() {
	listerner, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("cant make a tcp connection ")

	}
	//contanct accepting new connection
	for {
		con, err := listerner.Accept()
		if err != nil {
			log.Fatal("connection can't be esatblished!!")
		}
		req, err := request.RequestFromReader(con)
		if err != nil {
			log.Fatal("error while reading req", err.Error())
		}
		fmt.Printf("origin:%s\n , method:%s\n", req.RequestLine.RequestTarget, req.RequestLine.Method)

		con.Close()
	}
}
