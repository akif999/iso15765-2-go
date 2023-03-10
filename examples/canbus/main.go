package main

import (
	"log"

	"github.com/akif999/iso15765-2/examples/vcan"
)

func main() {
	// This program just behaving as gateway.
	var rxIDs map[uint16]struct{}
	for i := 0x000; i < 0x800; i++ {
		rxIDs[uint16(i)] = struct{}{}
	}
	type canMsg struct {
		canid uint16
		dlc   uint8
		data  []byte
	}
	ch := make(chan canMsg, 1)
	rxHandler := func(canid uint16, dlc uint8, data []byte) error {
		ch <- canMsg{canid, dlc, data}
		return nil
	}
	canbus := vcan.New(rxIDs, `\\.\pipe\canbus_out`, `\\.\pipe\canbus_in`, rxHandler)

	go func() {
		for {
			msg := <-ch
			canbus.Tx(msg.canid, msg.dlc, msg.data)
		}
	}()

	err := canbus.WaitForReception()
	if err != nil {
		log.Fatal(err)
	}
}
