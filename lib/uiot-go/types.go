package uiot

import (
	"fmt"
	"sync"

	proto "github.com/TrevorFarrelly/u-iot/proto"
)

// Device wraps the protobuf representation of a device's information
type Device struct {
	Name  string
	Type  Type
	Room  Room
	Funcs map[string]*Func

	// these fields are set on creation/receipt of a new device
	remote bool
	addr   string
	port   int
}

// create a local Device struct from the protobuf representation
func deviceFromProto(p *proto.DevInfo) *Device {
	// build device struct
	dev := &Device{
		Name:   p.Meta.Name,
		Type:   Type(p.Meta.Type),
		Room:   Room(p.Meta.Room),
		Funcs:  make(map[string]*Func),
		remote: true,
		addr:   p.Addr,
		port:   int(p.Port),
	}
	// add all functions to map
	for _, f := range p.Funcs {
		dev.Funcs[f.Name] = &Func{}
		// add all params to function
		for _, p := range f.Params {
			dev.Funcs[f.Name].Params = append(dev.Funcs[f.Name].Params, Param{int(p.Min), int(p.Max)})
		}
	}
	return dev
}

// add a function to this device
func (d *Device) AddFunction(name string, f func(...int), p ...Param) error {
	if d.remote {
		return fmt.Errorf("Cannot modify remote device")
	}
	d.Funcs[name] = &Func{f, p}
	return nil
}

// get the protobuf representation of this device
func (d Device) asProto() *proto.DevInfo {
	// construct device info
	pd := &proto.DevInfo{
		Port: uint32(d.port),
		Addr: d.addr,
		Meta: &proto.Meta{
			Type: uint32(d.Type),
			Room: uint32(d.Room),
			Name: d.Name,
		},
	}
	// add all functions to this device
	for name, f := range d.Funcs {
		pf := &proto.Func{
			Name: name,
		}
		// add all parameters to this function
		for _, p := range f.Params {
			pf.Params = append(pf.Params, &proto.Param{
				Min: uint32(p.Min),
				Max: uint32(p.Max),
			})
		}
		pd.Funcs = append(pd.Funcs, pf)
	}
	return pd
}

// returns the IP and port combo for this device
func (d Device) getFullAddress() string {
	return fmt.Sprintf("%s:%d", d.addr, d.port)
}

// print formatting
func (d Device) String() string {
	ret := fmt.Sprintf("(%s, %s) %s:", d.Room, d.Type, d.Name)
	for name, f := range d.Funcs {
		ret += fmt.Sprintf("  %s%s", name, f)
	}
	return ret
}

// Func wraps the protobuf representation of a device's function
type Func struct {
	F      func(...int)
	Params []Param
}

// print formatting
func (f Func) String() string {
	ret := fmt.Sprintf("( ")
	for _, p := range f.Params {
		ret += fmt.Sprintf("%s ", p)
	}
	ret += ")"
	return ret
}

// Param wraps the protobuf representation of a function's parameter
type Param struct {
	Min int
	Max int
}

// print formatting
func (p Param) String() string {
	return fmt.Sprintf("%d-%d", p.Min, p.Max)
}

// Network contains all of the known devices on the network, as well as an event-driven
// interface to detect when new devices connect
type Network struct {
	mux      sync.Mutex
	devs     []*Device
	eventing bool
	event    chan *Device
}

// Get all known remote devices
func (n Network) GetDevices() []*Device {
	return n.devs
}

// Add a new known device to the network
func (n *Network) addDevice(new *Device) error {
	n.mux.Lock()
	defer n.mux.Unlock()
	// determine if new device is already known by comparing its IP and port to other known devices
	for _, dev := range n.devs {
		if dev.getFullAddress() == new.getFullAddress() {
			return fmt.Errorf("%s (addr %s) already exists on the network", new.Name, new.getFullAddress())
		}
	}
	// add device to list
	n.devs = append(n.devs, new)
	// send event if it's enabled
	if n.eventing {
		n.event <- new
	}
	return nil
}

// enable the event interface
func (n *Network) EnableEvents() chan *Device {
	n.eventing = true
	return n.event
}

// Type represents the various types of devices that can exist on the network
type Type int

const (
	Light Type = iota
	Outlet
	Speaker
	Screen
	Controller
	OtherType
)

func (t Type) String() string {
	return [...]string{"Light", "Outlet", "Speaker", "Screen", "Controller", "Other"}[t]
}

// Room represents the various rooms a device can be placed in
type Room int

const (
	Living Room = iota
	Dining
	Bed
	Bath
	Kitchen
	Foyer
	Closet
	OtherRoom
)

func (r Room) String() string {
	return [...]string{"Living Room", "Dining Room", "Bedroom", "Bathroom", "Kitchen", "Foyer", "Closet", "Other"}[r]
}
