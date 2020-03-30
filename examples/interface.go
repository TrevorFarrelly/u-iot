package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	uiot "github.com/TrevorFarrelly/u-iot/lib/uiot-go"
)

func main() {
	name := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])

	// create a new device.
	// We specify the port the RPC server will listen on, as well as a device type
	// and location for convenience
	d := uiot.NewDevice(name, port, uiot.Controller, uiot.OtherRoom)

	// connect to the network.
	// we get a Network instance back. This instance is updated in the background
	// automaticaly as new devices connect.
	net, err := uiot.Bootstrap(d)
	if err != nil {
		log.Printf("could not bootstrap: %v", err)
	}

	// print network info after we are connected
	fmt.Printf("Starting u-iot UI...\nTriggering syntax: 'device function param1 param2 param3...'\n")
	time.Sleep(1 * time.Second)
	fmt.Printf("Available devices:\n")
	for _, d := range net.GetDevices() {
		fmt.Printf(" * %s\n", d)
	}
	fmt.Printf("> ")

	// get input from the user
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {

		// parse input
		input := strings.Split(s.Text(), " ")
		var params []int
		for _, param := range input[2:] {
			p, err := strconv.Atoi(param)
			if err != nil {
				log.Printf("Error parsing parameter %s: not an int\n", param)
				break
			}
			params = append(params, p)
		}

		// call function
		dev, err := net.GetDevice(input[0])
		if err != nil {
			log.Printf("Could not find device: %v", err)
		}
		dev.CallFunc(input[1], params...)

		fmt.Printf("> ")
	}
}
