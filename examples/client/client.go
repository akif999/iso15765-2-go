package main

import (
	"fmt"
	"log"

	"github.com/akif999/iso15765-2/examples/vcan"
)

func main() {
	rxIDs := map[uint16]struct{}{}
	rxIDs[0x6FF] = struct{}{}

	rxHandler := func(canid uint16, dlc uint8, data []byte) error {
		fmt.Printf("%03X %01X ", canid, dlc)
		for i, d := range data {
			if i == len(data) {
				fmt.Printf("%02X\n", d)
			} else {
				fmt.Printf("%02X ", d)
			}
		}
		return nil
	}
	client := vcan.New(rxIDs, `\\.\pipe\canbus_out`, `\\.\pipe\canbus_in`, rxHandler)

	// go func() {
	// 	for {
	// 		server.Tx(0x6FF, 8, []byte{0x01, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE})
	// 		time.Sleep(1000 * time.Millisecond)
	// 	}
	// }()

	err := client.WaitForReception()
	if err != nil {
		log.Fatal(err)
	}
}
