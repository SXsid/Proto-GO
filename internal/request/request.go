package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

// the hhtp mehtod consist
// 1) request line =>  method (get/post/put) rqust target (url/ endpoint) httpversiio (http1.1)
type Request struct {
	RequestLine RequestLine
	ParseStatus parseStatus
}
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}
type parseStatus int

const (
	initialized parseStatus = iota
	done
)
const maxBufferSize = 8 << 20

// constructor
func NewRequest() *Request {
	return &Request{
		ParseStatus: initialized,
	}
}

var Error_BAD_START_LINE = fmt.Errorf("somethign is wrong with start line")
var Error_EMPTY_REQUEST_LINE = fmt.Errorf("the request url is empty")
var ERROR_MALFORMED_START_LINE = fmt.Errorf("invalid start line / something is pararms are missing")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("inccorect http version")
var ERROR_UNSUPPORTED_HTTP_METHOD = fmt.Errorf("invalid method")
var ERROR_PARSING_IN_DONE_STATE = fmt.Errorf("trying to read data in a done state")
var ERROR_INVALID_STATE = fmt.Errorf("parser state is unknown")

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := NewRequest()
	//reading from stram of data

	buff := make([]byte, maxBufferSize)

	buffIndex := 0
	for req.ParseStatus != done {
		n, err := reader.Read(buff[buffIndex:])
		//what to do in read failure1
		if err != nil {
			return nil, err
		}
		buffIndex += n
		parsedIndex, err := req.parse(buff[:buffIndex])
		if err != nil {
			return nil, err
		}
		copy(buff, buff[parsedIndex:buffIndex])
		buffIndex -= parsedIndex
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

loop:
	for {
		switch r.ParseStatus {
		case done:
			break loop
		case initialized:
			reqLine, readlineindex, err := parseRequestLine(string(data[read:]))
			if err != nil {
				return 0, err
			}
			if readlineindex == 0 {
				break loop
			}
			r.RequestLine = reqLine
			read += readlineindex
			r.ParseStatus = done
		default:
			break loop
		}
	}
	return read, nil
}
func parseRequestLine(reqLine string) (RequestLine, int, error) {

	idx := strings.Index(reqLine, "\r\n")
	if idx == -1 {
		return RequestLine{}, 0, nil
	}
	startLine := reqLine[:idx]
	read := idx + len("\r\n") // accoutn for the limiters
	data := strings.Split(startLine, " ")
	if len(data) != 3 {
		return RequestLine{}, 0, ERROR_MALFORMED_START_LINE
	}
	if !isAllCaptilized(data[0]) {
		return RequestLine{}, 0, ERROR_UNSUPPORTED_HTTP_METHOD
	}
	if data[2] != "HTTP/1.1" {
		return RequestLine{}, 0, ERROR_UNSUPPORTED_HTTP_VERSION
	}

	return RequestLine{
		Method:        data[0],
		RequestTarget: data[1],
		HttpVersion:   strings.Split(data[2], "/")[1],
	}, read, nil
}

func isAllCaptilized(str string) bool {
	if len(str) == 0 {
		return false
	}

	for _, r := range str {
		if !unicode.IsLetter(r) || !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}
