package main

import (
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/TrevorFarrelly/rpi"
	"github.com/TrevorFarrelly/u-iot"
)

var (
	// RGB codes for a smooth rainbow in the cycle function
	cycle_vals = [360]int{0, 0, 0, 0, 0, 1, 1, 2, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13,
		15, 17, 18, 20, 22, 24, 26, 28, 30, 32, 35, 37, 39, 42, 44, 47, 49, 52, 55, 58,
		60, 63, 66, 69, 72, 75, 78, 81, 85, 88, 91, 94, 97, 101, 104, 107, 111, 114,
		117, 121, 124, 127, 131, 134, 137, 141, 144, 147, 150, 154, 157, 160, 163, 167,
		170, 173, 176, 179, 182, 185, 188, 191, 194, 197, 200, 202, 205, 208, 210, 213,
		215, 217, 220, 222, 224, 226, 229, 231, 232, 234, 236, 238, 239, 241, 242, 244,
		245, 246, 248, 249, 250, 251, 251, 252, 253, 253, 254, 254, 255, 255, 255, 255,
		255, 255, 255, 254, 254, 253, 253, 252, 251, 251, 250, 249, 248, 246, 245, 244,
		242, 241, 239, 238, 236, 234, 232, 231, 229, 226, 224, 222, 220, 217, 215, 213,
		210, 208, 205, 202, 200, 197, 194, 191, 188, 185, 182, 179, 176, 173, 170, 167,
		163, 160, 157, 154, 150, 147, 144, 141, 137, 134, 131, 127, 124, 121, 117, 114,
		111, 107, 104, 101, 97, 94, 91, 88, 85, 81, 78, 75, 72, 69, 66, 63, 60, 58, 55,
		52, 49, 47, 44, 42, 39, 37, 35, 32, 30, 28, 26, 24, 22, 20, 18, 17, 15, 13, 12,
		11, 9, 8, 7, 6, 5, 4, 3, 2, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// RGB pinout
	r = rpi.BoardToPin(36)
	g = rpi.BoardToPin(38)
	b = rpi.BoardToPin(40)

	// maximum PWM values for white balance. Normally done with resistors,
	// I did not have the appropriate ones on hand so I simulated it in software.
	rmax = 100
	gmax = 60
	bmax = 60

	// current PWM values
	rval = 100
	gval = 60
	bval = 60

	// determine if light is on or not
	lit = false
)

// Turn on the light, setting to the current color
func on(unused ...int) {
	lit = true
	log.Printf("turning light on...\n")
	rpi.PinMode(r, rpi.PWM_OUTPUT)
	rpi.PinMode(g, rpi.PWM_OUTPUT)
	rpi.PinMode(b, rpi.PWM_OUTPUT)
	rpi.DigitalWrite(r, rpi.HIGH, rval)
	rpi.DigitalWrite(g, rpi.HIGH, gval)
	rpi.DigitalWrite(b, rpi.HIGH, bval)
}

// turn off the light
func off(unused ...int) {
	lit = false
	log.Printf("turning light off...\n")
	rpi.PinMode(r, rpi.PWM_OUTPUT)
	rpi.PinMode(g, rpi.PWM_OUTPUT)
	rpi.PinMode(b, rpi.PWM_OUTPUT)
	rpi.DigitalWrite(r, rpi.LOW, 0)
	rpi.DigitalWrite(g, rpi.LOW, 0)
	rpi.DigitalWrite(b, rpi.LOW, 0)
}

// set the light to a new color
func color(vals ...int) {
	rval = int(math.Round(float64(vals[0]) / 255 * rmax))
	gval = int(math.Round(float64(vals[1]) / 255 * gmax))
	bval = int(math.Round(float64(vals[2]) / 255 * bmax))
	log.Printf("updating values to %d, %d, %d", rval, gval, bval)
	if lit {
		on()
	}
}

// cycle the LED in a rainbow, stopping when off() is called
func cycle(unused ...int) {
	go func() {
		lit = true
		log.Printf("cycling colors...\n")
		rpi.PinMode(r, rpi.PWM_OUTPUT)
		rpi.PinMode(g, rpi.PWM_OUTPUT)
		rpi.PinMode(b, rpi.PWM_OUTPUT)
		i := 0
		for {
			rpi.DigitalWrite(r, rpi.HIGH, cycle_vals[(i+120)%360])
			rpi.DigitalWrite(g, rpi.HIGH, cycle_vals[i%360])
			rpi.DigitalWrite(b, rpi.HIGH, cycle_vals[(i+240)%360])
			time.Sleep(20 * time.Millisecond)
			if !lit {
				break
			}
			i = i + 1
		}
	}()
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

	// default LED to off
	off()

	// create new u-iot device with all functions
	d := uiot.NewDevice(name, uiot.Light, uiot.Living)
	d.AddFunction("on", on)
	d.AddFunction("off", off)
	d.AddFunction("cycle", cycle)
	// color takes RGB values, i.e. 3 ints in the range 0-255
	d.AddFunction("color", color, uiot.Param{0, 255}, uiot.Param{0, 255}, uiot.Param{0, 255})

	// connect to the network
	_, err := uiot.Bootstrap(d, port)
	if err != nil {
		log.Printf("could not bootstrap: %v", err)
	}
}
