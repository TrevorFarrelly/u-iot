package main

import (
	"log"
	"os"
	"strconv"

	uiot "github.com/TrevorFarrelly/u-iot/lib/uiot-go"
)

// functions that this device performs. Due to Go's strict typing, parameters
// must be variadic. As long as the signature is defined properly, the library
// will verify that you get the number of variables you want, and they are within
// the range you want.
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
	name := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])

	// create a new device.
	// We specify the port the RPC server will listen on, as well as a device type
	// and location for convenience
	d := uiot.NewDevice(name, port, uiot.Light, uiot.Living)

	// add a function to the new device
	// We provide a name for the function and the expected parameters. Each
	// parameter has a range of values it can take. u-iot handles input sanitization
	// internally.
  d.AddFunction("hello0", Hello0)
  d.AddFunction("hello1", Hello1, uiot.Param{0, 256})
	d.AddFunction("hello3", Hello3, uiot.Param{0, 85}, uiot.Param{86, 171}, uiot.Param{172, 256})

	// connect to the network.
	// we get a Network instance back. This instance is updated in the background
	// automaticaly as new devices connect.
	net, err := uiot.Bootstrap(d)
	if err != nil {
		log.Printf("could not bootstrap: %v", err)
	}
	// print network info when a new device joins.
	// the Network instance contains a channel that triggers when a new device is
	// added to its internal list. This is an easy, if basic, way to create an
	// eventing interface.
	for dev := range net.EnableEvents() {
		log.Printf("New device detected: %s\n", dev.Name)
		log.Printf("Network:\n")
		for _, d := range net.GetDevices() {
			log.Printf(" * %s\n", d)
		}
	}
}
