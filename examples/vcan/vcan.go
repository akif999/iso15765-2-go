package vcan

import (
	"io/ioutil"
	"log"

	"github.com/Microsoft/go-winio"
)

type rxHandler func(canid uint16, dlc uint8, data []byte) error

type Node struct {
	rxIDs map[uint16]struct{}

	dialPipe   string
	listenPipe string

	txMsg []byte
	// rxMsg   []byte
	handler rxHandler
}

const ()

var (
	pipeConfig = winio.PipeConfig{
		SecurityDescriptor: "S:(ML;;NW;;;LW)D:(A;;0x12019f;;;WD)",
		InputBufferSize:    4096,
		OutputBufferSize:   4096,
	}
)

func New(rxIDs map[uint16]struct{}, dialPipe, listenPipe string, handler rxHandler) *Node {
	n := Node{
		rxIDs: rxIDs,

		dialPipe:   dialPipe,
		listenPipe: listenPipe,

		txMsg: make([]byte, 3+1+8), // canid + dlc + data
		// rxMsg: make([]byte, 3+1+8),
		handler: handler,
	}
	return &n
}

func (n *Node) Tx(canid uint16, dlc uint8, data []byte) error {
	n.txMsg[0] = byte((canid & 0x0700) >> 8)
	n.txMsg[1] = byte((canid & 0x00FF) >> 0)
	n.txMsg[2] = dlc
	copy(n.txMsg[3:8], data[:8])
	go func() {
		for {
			conn, err := winio.DialPipe(n.dialPipe, nil)
			if conn != nil && err == nil {
				defer conn.Close()
				conn.Write(n.txMsg)
				return
			}
		}
	}()
	return nil
}

func (n *Node) WaitForReception() error {
	listener, err := winio.ListenPipe(n.listenPipe, &pipeConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		// This is blocking API.
		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			return err
		}
		bytes, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}
		canid := uint16((uint16(bytes[0]&0x07) << 8) | (uint16(bytes[1]&0xFF) << 0))
		dlc := bytes[2]
		data := bytes[3:8]
		_, ok := n.rxIDs[canid]
		if ok {
			err = n.handler(canid, dlc, data)
		}
		if err != nil {
			return err
		}
	}
}
