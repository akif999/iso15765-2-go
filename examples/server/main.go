package main

import (
	"context"
	"fmt"
	"log"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/candevice"
	"go.einride.tech/can/pkg/socketcan"
)

func main() {
	d, err := candevice.New("can0")
	if err != nil {
		log.Fatal(err)
	}
	err = d.SetBitrate(250000)
	if err != nil {
		log.Fatal(err)
	}
	err = d.SetUp()
	if err != nil {
		log.Fatal(err)
	}
	defer d.SetDown()

	conn, err := socketcan.DialContext(context.Background(), "can", "can0")
	if err != nil {
		log.Fatal(err)
	}

	recv := socketcan.NewReceiver(conn)
	for recv.Receive() {
		frame := can.Frame{}
		fmt.Println(frame.String())
	}
}
