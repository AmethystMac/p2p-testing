package stream

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

const ChatProtocolId protocol.ID = "/chat/1.0.0"

func ChatProtocolHandler(streamRef network.Stream)  {

	readWriter := bufio.NewReadWriter(bufio.NewReader(streamRef), bufio.NewWriter(streamRef))

	fmt.Print("\nReceived connection from a peer. Type and press [Enter] to send.\n\n")
	
	go ReadStream(readWriter)
	go WriteStream(readWriter)
}

func ReadStream(readWriter *bufio.ReadWriter) {
	
	reader := readWriter.Reader

	for {
		contents, err := reader.ReadString('\n')
		if err != nil {
			log.Panicln(err)
		}

		fmt.Printf("%s> ", contents)
	}
}

func WriteStream(readWriter *bufio.ReadWriter) {

	cmdLine := bufio.NewReader(os.Stdin)
	writer := readWriter.Writer

	for {
		fmt.Print("> ")
		input, err := cmdLine.ReadString('\n')
		if err != nil {
			log.Panicln(err)
		}

		_, err = writer.WriteString(input)
		if err != nil {
			log.Panicln(err)
		}
		
		writer.Flush()
	}
}