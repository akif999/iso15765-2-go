package main

import (
	"context"
	"log"
	"time"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

// Run the following command before running this program.
// $ sudo modprobe can
// $ sudo modprobe can_raw
// $ sudo modprobe vcan
// $ sudo ip link add dev vcan0 type vcan
// $ sudo ip link set up vcan0
func main() {
	frame := can.Frame{
		ID:         0x546,
		Length:     8,
		Data:       can.Data{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF},
		IsRemote:   false,
		IsExtended: false,
	}

	conn, err := socketcan.DialContext(context.Background(), "can", "vcan0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for {
		tx := socketcan.NewTransmitter(conn)
		err = tx.TransmitFrame(context.Background(), frame)
		time.Sleep(1000 * time.Millisecond)
	}
}
