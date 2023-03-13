package main

import (
	"context"
	"fmt"
	"log"

	"go.einride.tech/can/pkg/socketcan"
)

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
