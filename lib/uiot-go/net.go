package uiot

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	proto "github.com/TrevorFarrelly/u-iot/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const (
	// Multicast information
	mcastaddr = "239.0.0.0:1024"
	mcastbuf  = 512
	mcastlen  = 2
)

// barebones information about a remote device. Used for async communication between
// the multicast server and RPC client.
type remote struct {
	addr string
	port int
}

// Handler for all incoming and outgoing multicasts
type mcastEndpoint struct {
	addr    *net.UDPAddr
	channel chan *remote
}

// receive incoming multicast messages
func (me *mcastEndpoint) recvMulticast() {
	// set up socket
	conn, err := net.ListenMulticastUDP("udp4", nil, me.addr)
	if err != nil {
		log.Fatalf("Could not start multicast server: %v", err)
	}
	conn.SetReadBuffer(mcastbuf)
	buf := make([]byte, mcastlen)
	// recieve packets
	for {
		_, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		// parse the data received and send it to the RPC client
		addr := strings.Split(src.String(), ":")[0]
		port := binary.BigEndian.Uint16(buf)
		me.channel <- &remote{addr, int(port)}
	}
}

// send our port to the multicast address
func (me *mcastEndpoint) sendMulticast(rpcPort int) error {
	// set up socket
	conn, err := net.DialUDP("udp4", nil, me.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// format and send data
	buf := make([]byte, mcastlen)
	binary.BigEndian.PutUint16(buf, uint16(rpcPort))
	conn.Write(buf)
	return nil
}

// start the multicast server and send our multicast message a moment later
func (me *mcastEndpoint) startMulticastService(rpcPort int) error {
	var err error
	me.addr, err = net.ResolveUDPAddr("udp4", mcastaddr)
	if err != nil {
		return err
	}

	go me.recvMulticast()
	time.Sleep(100 * time.Millisecond)
	if err := me.sendMulticast(rpcPort); err != nil {
		return err
	}

	return nil
}

// Handler for all incoming and outgoing RPCs
type rpcEndpoint struct {
	proto.UnimplementedDeviceServer
	local   *Device
	network *Network
	channel chan *remote
}

// RPC implementations

// Bootstrap a remote device. Add it to our network, then send our device info back
func (re *rpcEndpoint) Bootstrap(ctx context.Context, remote *proto.DevInfo) (*proto.DevInfo, error) {
	// parse remote device information
	device := deviceFromProto(remote)

	// get address and port info from the context
	p, ok := peer.FromContext(ctx)
	if ok {
		addr := strings.Split(p.Addr.String(), ":")
		device.addr = addr[0]
	} else {
		return nil, fmt.Errorf("Could not parse remote device information")
	}

	// add remote device info to our Network
	return re.local.asProto(), re.network.addDevice(device)
}

// Call a local function, triggered by a remote device
func (re *rpcEndpoint) CallFunc(ctx context.Context, funcinfo *proto.FuncCall) (*proto.FuncRet, error) {
	// get function from local device
	f, ok := re.local.Funcs[funcinfo.Name]
	if !ok {
		return &proto.FuncRet{}, fmt.Errorf("Device does not have requested function %s", funcinfo.Name)
	}
	// parse parameters
	var params []int
	for _, p := range funcinfo.Params {
		params = append(params, int(p))
	}
	// call function
	f.F(params...)
	return &proto.FuncRet{}, nil
}
}

// start the RPC server
func (re *rpcEndpoint) listenRPC(port int) error {
	// set up socket
	sock, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	// start RPC server
	server := grpc.NewServer()
	proto.RegisterDeviceServer(server, re)
	server.Serve(sock)
	return nil
}

// send Bootstrap RPC to devices we receive multicasts from
func (re *rpcEndpoint) sendBootstrap() {
	// get remote info from multicast server
	for r := range re.channel {
		// connect to remote device
		addr := fmt.Sprintf("%s:%d", r.addr, r.port)
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			continue
		}
		client := proto.NewDeviceClient(conn)
		ctx := context.Background()
		// send bootstrap info
		remote, err := client.Bootstrap(ctx, re.local.asProto())
		if err != nil {
			continue
		}
		// add remote device info to our Network
		device := deviceFromProto(remote)
		device.addr = r.addr
		device.port = r.port
		device.remote = true
		re.network.addDevice(device)
	}
}

// set up the RPC server and prepare to send Bootstrap RPCs when we receive a multicast
func (re *rpcEndpoint) startRPCService(rpcPort int) error {
	go re.listenRPC(rpcPort)
	time.Sleep(100 * time.Millisecond)
	go re.sendBootstrap()
	return nil
}
