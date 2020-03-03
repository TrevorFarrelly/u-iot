// Experimental program demonstrating the future bootstrapping process that will
// be used by u-iot.

package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// Future library code

// multicast information
const (
	mcastaddr = "239.0.0.0:1024"
	mcastlen  = 512
)

// local device state
var (
	me = &Device{}
)

// RPC server struct and methods
type uiotServer struct {
	UnimplementedDeviceServer
	mux  *sync.Mutex
	devs []*DevInfo
}

// receive a bootstrap RPC from a remote device
func (s *uiotServer) Bootstrap(ctx context.Context, dev *DevInfo) (*DevInfo, error) {
	log.Printf("Recv'd device info from GRPC client")
	p, ok := peer.FromContext(ctx)
	if ok {
		addr := strings.Split(p.Addr.String(), ":")
		portstr, _ := strconv.Atoi(addr[1])
		dev.Id.Address = addr[0]
		dev.Id.Port = uint32(portstr)
	} else {
		log.Printf("Could not get client info from context")
	}
	s.addDevice(dev)
	s.showDevices()
	return ProtoFromDevice(me), nil
}

// add a device to our list of known devices, in a thread-safe manner
func (s *uiotServer) addDevice(new *DevInfo) {
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, dev := range s.devs {
		if dev.Name == new.Name && dev.Id.Address == new.Id.Address {
			return
		}
	}
	s.devs = append(s.devs, new)
}

// show all devices that we know about
func (s *uiotServer) showDevices() {
	fmt.Printf("Current devices:\n")
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, dev := range s.devs {
		fmt.Printf("%s - ", dev.Name)
		for _, f := range dev.Funcs {
			fmt.Printf("%s(", f.Name)
			for i, p := range f.Params {
				fmt.Printf("%d-%d", p.Min, p.Max)
				if i < len(f.Params)-1 {
					fmt.Printf(", ")
				}
			}
			fmt.Printf(") ")
		}
		fmt.Printf("\n")
	}
}

// abstractions of the protobuf structs: Devices contain 0+ Funcs, Funcs contain 0+ Params
type Param struct {
	min uint32
	max uint32
}
type Func struct {
	Name   string
	F      func(...int)
	Params []Param
}
type Device struct {
	Name  string
	Funcs []Func
}

// IP information about a remote device
type Remote struct {
	ip   string
	port int
}

// save our user's device info to local state
func Register(devname string, funcs ...Func) {
	me.Name = devname
	me.Funcs = funcs
}

// build a protobuf DevInfo from the provided device
func ProtoFromDevice(dev *Device) *DevInfo {
	// build protobuf DevInfo
	rpc := DevInfo{
		Id: &ID{
			Id: 0,
		},
		Name: dev.Name,
	}
	// add all funcs
	for i, f := range dev.Funcs {
		new := &FuncDef{
			Id:   uint32(i),
			Name: f.Name,
		}
		// add all parameters for this func
		for _, p := range f.Params {
			new.Params = append(new.Params, &ParamDef{
				Min: p.min,
				Max: p.max,
			})
		}
		rpc.Funcs = append(rpc.Funcs, new)
	}

	return &rpc
}

// build a Device from the provided protobuf DevInfo
func DeviceFromProto(rpc *DevInfo) *Device {
	// initialize device struct
	dev := Device{Name: rpc.Name}
	// add all funcs
	for _, f := range rpc.Funcs {
		new := Func{Name: f.Name}
		// add all parameters for this func
		for _, p := range f.Params {
			new.Params = append(new.Params, Param{
				min: p.Min,
				max: p.Max,
			})
		}
		dev.Funcs = append(dev.Funcs, new)
	}
	return &dev
}

// receive multicasts from remote devices
func recvMulticast(addr *net.UDPAddr, remote chan Remote) {
	// set this socket to listen for multicasts on the specified address
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatalf("Could not join multicast group: %s\n", err)
	}
	conn.SetReadBuffer(mcastlen)

	// receive any incoming multicast messages
	for {
		buf := make([]byte, 2)
		_, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Failed to read message: %s\n", err)
		}
		// send results on channel
		ip := strings.Split(src.String(), ":")[0]
		port := binary.BigEndian.Uint16(buf)
		log.Printf("Recv'd GRPC addr %s:%d from multicast addr %s", ip, port, mcastaddr)
		remote <- Remote{ip, int(port)}
	}
}

// send multicasts to remote devices
func sendMulticast(addr *net.UDPAddr, rpcport int) {

	// set this socket to send on the multicast address
	conn, err := net.DialUDP("udp4", nil, addr)
	defer conn.Close()
	if err != nil {
		log.Fatalf("Could not set up socket: %s\n", conn)
	}

	// send the port we are receiving RPC's on
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(rpcport))
	conn.Write(b)
	log.Printf("Sent GRPC port %d to multicast addr %s", rpcport, mcastaddr)
}

// start RPC server to listen for incoming messages
func recvBootstrapRPC(serv *uiotServer, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to start RPC server: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterDeviceServer(grpcServer, serv)
	grpcServer.Serve(lis)
}

// send RPCs to devices that we receive a multicast from
func sendBootstrapRPC(serv *uiotServer, remote chan Remote) {
	// get server info from multicast
	for r := range remote {
		// connect to server
		server := fmt.Sprintf("%s:%d", r.ip, r.port)
		conn, err := grpc.Dial(server, []grpc.DialOption{grpc.WithInsecure()}...)
		if err != nil {
			log.Printf("Could not dial %s: %s\n", server, err)
		}
		client := NewDeviceClient(conn)
		ctx := context.Background()
		// request device information
		device, err := client.Bootstrap(ctx, ProtoFromDevice(me))
		log.Printf("Recv'd device info from GRPC server %s:%d", r.ip, r.port)
		if err != nil {
			log.Printf("%v.Bootstrap() failed: %s", client, err)
		}

		// add remote info to device
		device.Id.Address = r.ip
		device.Id.Port = uint32(r.port)
		// save device to internal database
		serv.addDevice(device)
		serv.showDevices()
	}
}

// Connect to other u-iot devices on the LAN
func Bootstrap(port int) *uiotServer {
	// set up UDP and RPC server
	addr, err := net.ResolveUDPAddr("udp4", mcastaddr)
	if err != nil {
		log.Fatalf("Could not resolve multicast addr: %s\n", err)
	}
	serv := &uiotServer{
		mux: &sync.Mutex{},
		devs: []*DevInfo{},
	}
	remote := make(chan Remote)
	// start servers
	go recvMulticast(addr, remote)
	go recvBootstrapRPC(serv, port)

	// wait for servers to spin up
	time.Sleep(100 * time.Millisecond)

	// send our bootstrapping messages
	go sendMulticast(addr, port)
	go sendBootstrapRPC(serv, remote)

	return serv
}

// Example program code Interface is subject to chane

var (
	name = flag.String("name", "", "device name")
	port = flag.Int("port", 2048, "port to receive RPCs on")
)

// functions that this device performs. Due to Go's strict typing, parameters
// must be variadic. As long as the signature is defined properly, the library
// will verify that you get the number of variables you want, and they are within
// the range you want.
func f0(args ...int) {
	log.Println("zero-arg func")
}
func f1(args ...int) {
	log.Println("one-arg func: %s", args[0])
}
func f3(args ...int) {
	log.Println("three-arg func: %s, %s, %s", args[0], args[1], args[2])
}

// define structs for each function. As stated above, as long as the definitions
// here match what is desired in the function, no extra verification will have to
// be done by the user.
func getFuncDefs() []Func {
	return []Func{
		{
			Name: "NoArgs",
			F:    f0,
		},
		{
			Name:   "OneArg",
			F:      f1,
			Params: []Param{{0, 255}},
		},
		{
			Name:   "ThreeArgs",
			F:      f3,
			Params: []Param{{0, 85}, {86, 170}, {171, 255}},
		},
	}
}

func main() {
	flag.Parse()
	Register(*name, getFuncDefs()...)
	Bootstrap(*port)
	// busy wait, implement any local device logic here (UIs, etc)
	for {}
}
