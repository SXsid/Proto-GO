package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// we are reading data slowly from buffer to mimc data is till comng and as we get one whole data packe we put it into a channel
func getLinesChannel(file io.ReadCloser) <-chan string {
	lines := make(chan string, 1)
	go func() {
		defer close(lines)
		defer file.Close()
		str := ""
		for {
			//buffer:holding place
			data := make([]byte, 8)
			//populating buffer while reutrn number of valid byte we can extract (error rep EOF)
			n, err := file.Read(data)
			if err != nil {
				break
			}
			tempdata := string(data[:n])
			parts := strings.Split(tempdata, "\r\n")
			//get all parts sep by spliter
			for i := 0; i < len(parts)-1; i++ {
				lines <- str + parts[i]
				str = ""
			}

			//last part of cuurent line
			str += parts[len(parts)-1]
		}
		//if after all read last chunk remains
		if str != "" {
			lines <- str
		}

	}()
	return lines
}

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
		for data := range getLinesChannel(con) {
			fmt.Printf("read:%s\n", data)
		}
		con.Close()
	}
}
