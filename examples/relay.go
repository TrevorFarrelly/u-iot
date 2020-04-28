package main

import (
	"log"
	"os"
	"strconv"

	"github.com/TrevorFarrelly/rpi"
	"github.com/TrevorFarrelly/u-iot"
)

var ctrl = rpi.BoardToPin(5)

// Set the control pin to HIGH.
// Arguments are required to work with u-iot, due to Go's strict type system.
// They are unused in this example
func on(unused ...int) {
	rpi.PinMode(ctrl, rpi.OUTPUT)
	rpi.DigitalWrite(ctrl, rpi.HIGH, -1)
}

// set the control pin to LOW
func off(unused ...int) {
	rpi.PinMode(ctrl, rpi.OUTPUT)
	rpi.DigitalWrite(ctrl, rpi.LOW, -1)
}

func main() {
	// Parse command-line input
	if len(os.Args) != 3 {
		log.Printf("Invalid arguments.\nUsage: %s name port\n", os.Args[0])
	}
	name := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Printf("Could not parse port number: %v\n", err)
	}

	// set up wiring library
	if err := rpi.WiringPiSetup(); err != nil {
		log.Printf("Could not setup GPIO: %v\n", err)
		return
	}

	// default pin to LOW
	off()

	// create new u-iot device with on/off functions
	d := uiot.NewDevice(name, uiot.Light, uiot.Living)
	d.AddFunction("on", on)
	d.AddFunction("off", off)

	// connect to the network
	_, err := uiot.Bootstrap(d, port)
	if err != nil {
		log.Printf("could not bootstrap: %v", err)
	}
}
