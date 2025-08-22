package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	parseStatus parseStatus
}

type RequestLine struct {
	Method string

	HttpVersion   string
	RequestTarget string
}

// constructor
func NewRequest() *Request {
	return &Request{
		parseStatus: parseInit,
	}
}

type parseStatus string

const (
	parseInit parseStatus = "init"
	parseDone parseStatus = "done"
)
const MaX_BUFF_SIZE = 8 << 20

var DELIMITOR = []byte("\r\n")
var ERROR_IVALID_STARTLINE = fmt.Errorf("incomplete/error full  start line ")
var ERROR_INVALID_PARSESTATUS = fmt.Errorf("unkown parse status")
var ERROR_UNSUPOORTED_HTTPVERSION = fmt.Errorf("unsupported http version ")
var ERROR_UNSUPOORTED_METHOD = fmt.Errorf("wrong methods")

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, MaX_BUFF_SIZE)
	bufflength := 0
	request := NewRequest()
	for request.parseStatus != parseDone {
		//read from the stream
		n, err := reader.Read(buff[bufflength:])
		if err != nil {
			return nil, err
		}
		bufflength += n
		read, err := request.parser(buff[:bufflength])
		if err != nil {
			return nil, err
		}
		copy(buff, buff[read:bufflength])
		bufflength -= read
	}
	return request, nil
}
func (r *Request) parser(data []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.parseStatus {
		case parseDone:
			break outer
		case parseInit:
			reqLine, readIndex, err := parseRequesLine(data)
			if err != nil {
				return 0, err
			}
			if readIndex == 0 {
				break outer
			}
			r.RequestLine = *reqLine
			r.parseStatus = parseDone
			read += readIndex
		default:
			return 0, ERROR_INVALID_PARSESTATUS
		}
	}
	return read, nil
}

func parseRequesLine(data []byte) (*RequestLine, int, error) {

	idx := bytes.Index(data, DELIMITOR)
	if idx == -1 {
		//continue reading
		return nil, 0, nil
	}
	parts := bytes.Split(data[:idx], []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_IVALID_STARTLINE
	}
	read := idx + len(DELIMITOR)
	//error
	if !isAllCapital(parts[0]) {
		return nil, 0, ERROR_UNSUPOORTED_METHOD
	}
	if string(parts[2]) != "HTTP/1.1" {
		return nil, 0, ERROR_UNSUPOORTED_HTTPVERSION
	}
	return &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   strings.Split(string(parts[2]), "/")[1],
	}, read, nil
}
func isAllCapital(method []byte) bool {
	if len(method) == 0 {
		return false
	}
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return false
		}
	}
	return true
}
