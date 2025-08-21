package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data          string
	currPost      int
	bytePerPacket int
}

// we will make this ifer the io.reader interface so we can pass donw our function
func (cr *chunkReader) Read(buf []byte) (howMuchPopulated int, err error) {
	if cr.currPost >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.currPost + cr.bytePerPacket
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n := copy(buf, cr.data[cr.currPost:endIndex])
	cr.currPost += n
	return n, nil
}

func TestRequestReader(t *testing.T) {
	//home patht testing
	reader := &chunkReader{
		data:          "GET / HTTP/1.1\r\nHOST: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		bytePerPacket: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	reader = &chunkReader{
		data:          "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		bytePerPacket: 1,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

}
