package main

import (
	"log"
	"time"

	"github.com/akif999/iso15765-2/examples/vcan"
)

func main() {
	var rxIDs map[uint16]struct{}
	rxIDs[0x7FF] = struct{}{}

	rxHandler := func(canid uint16, dlc uint8, data []byte) error {
		return nil
	}
	server := vcan.New(rxIDs, `\\.\pipe\canbus_out`, `\\.\pipe\canbus_in`, rxHandler)

	go func() {
		for {
			server.Tx(0x6FF, 8, []byte{0x01, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE})
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	err := server.WaitForReception()
	if err != nil {
		log.Fatal(err)
	}
}
