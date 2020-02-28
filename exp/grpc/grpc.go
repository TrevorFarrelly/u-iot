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
	"strings"
	"sync"
	"time"

	"github.com/TrevorFarrelly/u-iot/proto"
	"google.golang.org/grpc"
)

// Future library code

const (
	mcastaddr = "239.0.0.0:1024"
	mcastlen  = 512
)

var (
	me = &Device{}
)

// RPC server struct and methods
type uiotServer struct {
	uiot.UnimplementedDeviceServer
	mux  *sync.Mutex
	devs []*uiot.DevInfo
}

func (s *uiotServer) Bootstrap(ctx context.Context, dev *uiot.DevInfo) (*uiot.DevInfo, error) {
	s.addDevice(dev)
	return ProtoFromDevice(me), nil
}

// add a device to the list of known devices
func (s *uiotServer) addDevice(new *uiot.DevInfo) {
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, dev := range s.devs {
		if dev.Id.Address == new.Id.Address {
			return
		}
	}
	s.devs = append(s.devs, new)
}

// get a device from the server
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

// abstractions of the protobuf structs
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

// information about a remote device
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
func ProtoFromDevice(dev *Device) *uiot.DevInfo {
	// build protobuf message
	rpc := uiot.DevInfo{
		Id: &uiot.ID{
			Id: 0,
		},
		Name: dev.Name,
	}
	// add all funcs
	for i, f := range dev.Funcs {
		new := &uiot.FuncDef{
			Id:   uint32(i),
			Name: f.Name,
		}
		// add all parameters for this func
		for _, p := range f.Params {
			new.Params = append(new.Params, &uiot.ParamDef{
				Min: p.min,
				Max: p.max,
			})
		}
		rpc.Funcs = append(rpc.Funcs, new)
	}

	return &rpc
}

// build a Device from the provided protobuf DevInfo
func DeviceFromProto(rpc *uiot.DevInfo) *Device {
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
}

// receive RPCs from devices responding to our multicast
func recvBootstrapRPC(serv *uiotServer, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to start RPC server: %v", err)
	}
	grpcServer := grpc.NewServer()
	uiot.RegisterDeviceServer(grpcServer, serv)
	grpcServer.Serve(lis)
}

// send RPCs to devices that we receive a multicast from
func sendBootstrapRPC(serv *uiotServer, remote chan Remote) {
	for {
		// get server info from multicast
		r := <-remote
		// connect to server
		server := fmt.Sprintf("%s:%d", r.ip, r.port)
		conn, err := grpc.Dial(server, []grpc.DialOption{grpc.WithInsecure()}...)
		if err != nil {
			log.Printf("Could not dial %s: %s\n", server, err)
		}
		client := uiot.NewDeviceClient(conn)
		ctx := context.Background()
		// request device
		device, err := client.Bootstrap(ctx, ProtoFromDevice(me))
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
func Bootstrap(port int) {
	// set up servers to respond to other devices
	addr, err := net.ResolveUDPAddr("udp4", mcastaddr)
	if err != nil {
		log.Fatalf("Could not resolve multicast addr: %s\n", err)
	}
	serv := &uiotServer{
		mux: &sync.Mutex{},
		devs: []*uiot.DevInfo{},
	}
	remote := make(chan Remote)
	go recvMulticast(addr, remote)
	go recvBootstrapRPC(serv, port)

	time.Sleep(100 * time.Millisecond)
	// send our bootstrapping messages
	go sendMulticast(addr, port)
	sendBootstrapRPC(serv, remote)
}

// Example program code

var (
	name = flag.String("name", "", "device name")
	port = flag.Int("port", 2048, "port to receive RPCs on")
)

// functions that this device performs
func f0(args ...int) {
	log.Println("zero-arg func")
}
func f1(args ...int) {
	log.Println("one-arg func: %s", args[0])
}
func f3(args ...int) {
	log.Println("three-arg func: %s, %s, %s", args[0], args[1], args[2])
}

// define structs for each function
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
}
