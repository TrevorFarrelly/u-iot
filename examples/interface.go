package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TrevorFarrelly/u-iot/lib/uiot-go"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Invalid arguments. Usage:\n%s name port\n", os.Args[0])
		return
	}
	name := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Invalid port number: %s\n", os.Args[2])
	}

	// create a new device with no functions to represent our interface
	d := uiot.NewDevice(name, port, uiot.Controller, uiot.OtherRoom)

	// connect to the network.
	net, err := uiot.Bootstrap(d)
	if err != nil {
		log.Printf("could not bootstrap: %v", err)
	}

	fmt.Printf("Starting u-iot UI...\nTriggering syntax:\n > CallOne device function param1 param2 param3...\n > CallAll type room function param1 param2 param3...\n\n")
	// print network info after we are connected
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

		// calling a single device
		if input[0] == "CallOne" {
			if len(input) < 3 {
				log.Printf("Not enough arguments to CallOne command.\n")
				fmt.Printf("> ")
				continue
			}
			var params []int
			// parse parameters
			for _, param := range input[3:] {
				p, err := strconv.Atoi(param)
				if err != nil {
					log.Printf("Error parsing parameter %s: not an int\n", param)
					break
				}
				params = append(params, p)
			}
			// call function
			dev, err := net.GetDevice(input[1])
			if err != nil {
				log.Printf("Could not find device: %v", err)
			}
			if err = dev.CallFunc(input[2], params...); err != nil {
				log.Printf("Could not call function: %v", err)
			}

			// calling multiple devices
		} else if input[0] == "CallAll" {
			if len(input) < 4 {
				log.Printf("Not enough arguments to CallAll command.\n")
				fmt.Printf("> ")
				continue
			}
			// parse provided type
			t, err := uiot.TypeFromString(input[1])
			if err != nil {
				log.Printf("%v", err)
				fmt.Printf("> ")
				continue
			}
			// parse provided room
			r, err := uiot.RoomFromString(input[2])
			if err != nil {
				log.Printf("%v", err)
				fmt.Printf("> ")
				continue
			}
			var params []int
			// parse parameters
			for _, param := range input[4:] {
				p, err := strconv.Atoi(param)
				if err != nil {
					log.Printf("%v", err)
					break
				}
				params = append(params, p)
			}
			// call function
			if err = net.CallAll(r, t, input[3], params...); err != nil {
				log.Printf("Could not call all functions: %v", err)
			}
		}

		fmt.Printf("> ")
	}
}
