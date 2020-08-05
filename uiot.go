// Copyright (C) <year>  <name of author>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package uiot provides functionality for building your own smart home network.
package uiot

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	// net stores all devices we know about
	network *Network
	re      *rpcEndpoint
	me      *mcastEndpoint
)

// NewDevice creates a new device with the provided name, room, and type.
// This device is intended to be used locally, as a place to define functions
// that can be called remotely.
func NewDevice(name string, t Type, r Room) *Device {
	return &Device{
		Name:   name,
		Type:   t,
		Room:   r,
		Funcs:  make(map[string]*Func),
		remote: false,
	}
}

// Bootstrap starts networking services and broadcasts device information dev to
// the rest of the network.
func Bootstrap(dev *Device, port int) (*Network, error) {
	// initialize the network struct
	network = &Network{
		event: make(chan *Event),
	}
	// add RPC port info to device struct
  dev.port = port
	// create the channel used to send data between the multicast service and the RPC service
	c := make(chan *remote)
	// start the RPC service
	re = &rpcEndpoint{local: dev, network: network, channel: c}
	if err := re.startRPCService(port); err != nil {
		return nil, err
	}
	// start the multicast service
	me = &mcastEndpoint{channel: c}
	if err := me.startMulticastService(port); err != nil {
		return nil, err
	}
	return network, nil
}

// Close notifies all known remote devices that this device is disconnecting from
// the network.
func Close() {
	for _, dev := range network.devs {
		re.sendQuit(re.local, dev)
	}
}

// CloseHandler adds a SIGTERM handler to automatically disconnect from the network.
func CloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		Close()
		os.Exit(0)
	}()
}
