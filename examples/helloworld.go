package main

import (
	"log"
	"os"
	"strconv"

	"github.com/TrevorFarrelly/u-iot"
)

// functions that this device performs. Due to Go's strict typing, parameters
// must be variadic. As long as the signature is defined properly when registering
// the function, the library will verify that you get the number of variables you
// want, and they are within the range you want.

func Hello0(args ...int) {
	log.Printf("Hello, World!")
}

func Hello1(args ...int) {
	log.Printf("Hello, World! You sent %d", args[0])
}

func Hello3(args ...int) {
	log.Printf("Hello, World! You sent (%d, %d, %d)", args[0], args[1], args[2])
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

	// create a new device.
	// We specify the port the RPC server will listen on, as well as a device type
	// and location tags for convenience
	d := uiot.NewDevice(name, uiot.Light, uiot.Living)

	// add functions to the new device
	// We provide a name for the function and the expected parameters. Each
	// parameter has a range of values it can take. u-iot handles input sanitization
	// internally.
	d.AddFunction("hello0", Hello0)
	d.AddFunction("hello1", Hello1, uiot.Param{0, 256})
	d.AddFunction("hello3", Hello3, uiot.Param{0, 85}, uiot.Param{86, 171}, uiot.Param{172, 256})

	// connect to the network.
	// we get a Network instance back. This instance is updated in the background
	// automaticaly as new devices connect.
	net, err := uiot.Bootstrap(d, port)
	if err != nil {
		log.Printf("could not bootstrap: %v", err)
	}

	uiot.CloseHandler()

	// print network info when a new device joins.
	// the Network instance contains a channel that triggers when a new device is
	// added to its internal list. This is an easy, if basic, way to create an
	// eventing interface.
	for e := range net.EnableEvents() {
		log.Printf("Remote Device %sed: %s\n", e.Type.String(), e.Dev.Name)
		log.Printf("Network:\n")
		for _, d := range net.GetDevices() {
			log.Printf(" * %s\n", d)
		}
	}
}
