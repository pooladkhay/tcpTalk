package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gordonklaus/portaudio"
)

const serverAddr = "mamad.local:4004"

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	c, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	portaudio.Initialize()
	defer portaudio.Terminate()

	h, err := portaudio.DefaultHostApi()
	errCheck(err)
	p := portaudio.HighLatencyParameters(h.DefaultInputDevice, nil)
	p.Input.Channels = 1

	stream, err := portaudio.OpenStream(p, func(in, out []float32) {
		buf := new(bytes.Buffer)
		errCheck(binary.Write(buf, binary.BigEndian, in))
		_, err := c.Write(buf.Bytes())
		errCheck(err)
	})
	errCheck(err)
	defer stream.Close()

	errCheck(stream.Start())

	<-sig
	fmt.Println("\nExiting...")
	errCheck(stream.Stop())

}

func errCheck(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
