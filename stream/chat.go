package stream

import (
	"bufio"
	"log"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

const ProtocolId protocol.ID = "/chat/1.0.0"

func ChatProtocolHandler(streamRef network.Stream)  {

	readWriter := bufio.NewReadWriter(bufio.NewReader(streamRef), bufio.NewWriter(streamRef))
	
	go ReadStream(readWriter)
	go WriteStream(readWriter)
}

func ReadStream(readWriter *bufio.ReadWriter) {
	
	reader := readWriter.Reader

	contents, err := reader.ReadString('\n')
	if err != nil {
		log.Panicln(err)
	}

	log.Println(contents)
}

func WriteStream(readWriter *bufio.ReadWriter) {

	writer := readWriter.Writer

	_, err := writer.WriteString("Hello, People!\n")
	if err != nil {
		log.Panicln(err)
	}
	
	writer.Flush()
}