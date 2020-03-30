package uiot

import ()

var (
	// net stores all devices we know about
	network *Network
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
		event: make(chan *Device),
	}
	// create the channel used to send data between the multicast service and the RPC service
	c := make(chan *remote)
	// start the RPC service
	re := rpcEndpoint{local: dev, network: network, channel: c}
	if err := re.startRPCService(dev.port); err != nil {
		return nil, err
	}
	// start the multicast service
	me := mcastEndpoint{channel: c}
	if err := me.startMulticastService(dev.port); err != nil {
		return nil, err
	}
	return network, nil
}
