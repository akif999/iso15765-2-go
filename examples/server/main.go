package main

import (
	"context"
	"fmt"
	"log"

	"go.einride.tech/can/pkg/socketcan"
)

// Run the following command before running this program.
// $ sudo modprobe can
// $ sudo modprobe can_raw
// $ sudo modprobe vcan
// $ sudo ip link add dev vcan0 type vcan
// $ sudo ip link set up vcan0
func main() {
	conn, err := socketcan.DialContext(context.Background(), "can", "vcan0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	recv := socketcan.NewReceiver(conn)
	for recv.Receive() {
		frame := recv.Frame()
		fmt.Println(frame.String())
	}
}
