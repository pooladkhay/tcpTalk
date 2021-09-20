package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gordonklaus/portaudio"
)

const port = "4004"

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		os.Exit(0)
	}()

	PORT := ":" + port
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		log.Fatalln(err)
	}

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Client connected: %s\n", c.RemoteAddr().String())

	portaudio.Initialize()
	defer portaudio.Terminate()

	h, err := portaudio.DefaultHostApi()
	if err != nil {
		log.Fatalln(err)
	}

	p := portaudio.HighLatencyParameters(nil, h.DefaultOutputDevice)
	p.Output.Channels = 1

	stream, err := portaudio.OpenStream(p, func(in, out []float32) {
		err := binary.Read(c, binary.BigEndian, out)
		if err != nil {
			os.Exit(1)
		}
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
