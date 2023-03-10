package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/Microsoft/go-winio"
)

const (
	nodeApipe string = `\\.\pipe\nodeA`
	nodeBpipe string = `\\.\pipe\nodeB`
)

func main() {
	var argNode string
	var msg string
	flag.StringVar(&argNode, "node", "nodeA", "choice of node")
	flag.StringVar(&msg, "msg", "Hello, world!!", "transmitting msg")
	flag.Parse()

	const (
		nodeA int = iota
		nodeB
	)
	var listenPipe string
	var dialPipe string
	listenPipe = nodeApipe
	dialPipe = nodeBpipe
	if argNode != "nodeA" {
		listenPipe = nodeBpipe
		dialPipe = nodeApipe
	}

	pipeConfig := winio.PipeConfig{
		SecurityDescriptor: "S:(ML;;NW;;;LW)D:(A;;0x12019f;;;WD)",
		InputBufferSize:    4096,
		OutputBufferSize:   4096,
	}

	fmt.Printf("listen: %s\nDial: %s\n", listenPipe, dialPipe)

	go func() {
		cnt := 0
		for {
			go func() {
				conn, err := winio.DialPipe(dialPipe, nil)
				if conn != nil && err == nil {
					defer conn.Close()
					conn.Write([]byte(fmt.Sprintf("%s: %d", msg, cnt)))
					cnt++
				}
			}()
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	listener, err := winio.ListenPipe(listenPipe, &pipeConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		// This is blocking API.
		conn, _ := listener.Accept()
		defer conn.Close()

		bytes, _ := ioutil.ReadAll(conn)
		fmt.Println(string(bytes))
	}
}
