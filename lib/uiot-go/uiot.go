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

// create a new device
func NewDevice(name string, rpcPort int, t Type, r Room) *Device {
	return &Device{
		Name:   name,
		Type:   t,
		Room:   r,
		Funcs:  make(map[string]*Func),
		port:   rpcPort,
		remote: false,
	}
}

// connect to the network
func Bootstrap(dev *Device) (*Network, error) {
	// initialize the network struct
	network = &Network{
		event: make(chan *Event),
	}
	// create the channel used to send data between the multicast service and the RPC service
	c := make(chan *remote)
	// start the RPC service
	re = &rpcEndpoint{local: dev, network: network, channel: c}
	if err := re.startRPCService(dev.port); err != nil {
		return nil, err
	}
	// start the multicast service
	me = &mcastEndpoint{channel: c}
	if err := me.startMulticastService(dev.port); err != nil {
		return nil, err
	}
	return network, nil
}

// disconnect from the network
func Close() {
	for _, dev := range network.devs {
		re.sendQuit(re.local, dev)
	}
}

// automatically disconnect when Ctrl-C is pressed
func CloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		Close()
		os.Exit(0)
	}()
}
