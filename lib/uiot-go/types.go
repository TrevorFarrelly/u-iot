package uiot

import (
	"context"
	"fmt"
	"sync"

	proto "github.com/TrevorFarrelly/u-iot/proto"
	"google.golang.org/grpc"
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

// call a device's function
func (d Device) CallFunc(name string, p ...int) error {
	// get function from device
	f, ok := d.Funcs[name]
	if !ok {
		return fmt.Errorf("device %s does not have function %s", d.Name, name)
	}

	// check parameters
	if len(p) != len(f.Params) {
		return fmt.Errorf("%s expects %d parameters, but %d were provided", name, len(f.Params), len(p))
	}
	var convP []uint32
	for i, param := range f.Params {
		if p[i] < param.Min || p[i] > param.Max {
			return fmt.Errorf("%d is out of range for parameter %d: %d-%d", p[i], i, param.Min, param.Max)
		}
		convP = append(convP, uint32(p[i]))
	}

	// call remote function
	addr := fmt.Sprintf("%s:%d", d.addr, d.port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := proto.NewDeviceClient(conn)
	ctx := context.Background()
	_, err = client.CallFunc(ctx, &proto.FuncCall{Name: name, Params: convP})
	if err != nil {
		return err
	}
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
	event    chan *Event
}

// get a device from the network
func (n Network) GetDevice(name string) (*Device, error) {
	n.mux.Lock()
	defer n.mux.Unlock()
	for _, dev := range n.devs {
		if dev.Name == name {
			return dev, nil
		}
	}
	return nil, fmt.Errorf("%s was not found on the network", name)
}

// Get all known remote devices
func (n Network) GetDevices() []*Device {
	return n.devs
}

// Get all devices that match the specified type and/or room
func (n Network) GetMatching(r Room, t Type) []*Device {
	// initialize return array
	ret := []*Device{}
	// iterate over all devices
	for _, d := range n.GetDevices() {
		t_match := true
		r_match := true
		// if type is set and is not a match, skip it
		if t != -1 && d.Type != t {
			t_match = false
		}
		// if room is set and is not a match, skip it
		if r != -1 && d.Room != r {
			r_match = false
		}
		// if both type and room match, add it to return list
		if t_match && r_match {
			ret = append(ret, d)
		}
	}
	return ret
}

// call a function on all matching devices
func (n Network) CallAll(r Room, t Type, name string, p ...int) error {
	// initialize error array
	errs := []error{}
	// iterate over all matching devices
	for _, d := range n.GetMatching(r, t) {
		// call function, adding error to array if we get one
		if err := d.CallFunc(name, p...); err != nil {
			errs = append(errs, err)
		}
	}
	// if we encountered any errors, return one containing all of them
	if len(errs) != 0 {
		ret := fmt.Sprintf("could not call '%s' on %d device(s):", name, len(errs))
		for _, err := range errs {
			ret += fmt.Sprintf("\n  %v", err)
		}
		return fmt.Errorf(ret)
	}
	return nil
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
		n.event <- &Event{Connect, new}
	}
	return nil
}

// Remove a device from the network
func (n *Network) removeDevice(old *Device) error {
	n.mux.Lock()
	defer n.mux.Unlock()
	// remove the device from the network if we have it
	for i, dev := range n.devs {
		if dev.getFullAddress() == old.getFullAddress() {
			// swap device to remove with the device at the end of the list
			n.devs[len(n.devs)-1], n.devs[i] = n.devs[i], n.devs[len(n.devs)-1]
			n.devs = n.devs[:len(n.devs)-1]
			// send event if it's enabled
			if n.eventing {
				n.event <- &Event{Disconnect, old}
			}
			return nil
		}
	}
	return fmt.Errorf("%s (addr %s) does not exist on the network", old.Name, old.getFullAddress())
}

// enable the event interface
func (n *Network) EnableEvents() chan *Event {
	n.eventing = true
	return n.event
}

// Event represents a change in state in the network. Used to push notifications
// to the user when they take advantage of the eventing interface.
type Event struct {
	Type EventType
	Dev  *Device
}

// EventType represents the different types of events that are supported
type EventType int

const (
	Connect EventType = iota
	Disconnect
)

func (t EventType) String() string {
	return [...]string{"Connect", "Disconnect"}[t]
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

func TypeFromString(s1 string) (Type, error) {
	if s1 == "*" {
		return -1, nil
	}
	for i, s2 := range [...]string{"Light", "Outlet", "Speaker", "Screen", "Controller", "Other"} {
		if s1 == s2 {
			return Type(i), nil
		}
	}
	return -1, fmt.Errorf("%s is not a valid device type", s1)
}

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

func RoomFromString(s1 string) (Room, error) {
	if s1 == "*" {
		return -1, nil
	}
	for i, s2 := range [...]string{"LivingRoom", "DiningRoom", "Bedroom", "Bathroom", "Kitchen", "Foyer", "Closet", "Other"} {
		if s1 == s2 {
			return Room(i), nil
		}
	}
	return -1, fmt.Errorf("%s is not a valid device room", s1)
}

func (r Room) String() string {
	return [...]string{"LivingRoom", "DiningRoom", "Bedroom", "Bathroom", "Kitchen", "Foyer", "Closet", "Other"}[r]
}
