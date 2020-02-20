// This program contains an experiment with Go multicasting syntax. Something
// similar will be used in the u-iot bootstrapping process

package main

import (
  "bufio"
  "flag"
  "fmt"
  "log"
  "net"
  "os"
)

const (
  saddr = "239.0.0.0:9999"
  maxlen = 8192
)

var (
  name = flag.String("name", "", "device name")
  port = flag.Int("port", 1024, "port number")
)

func getLocalIP() net.IP {
  conn, err := net.Dial("udp", "1.1.1.1:80")
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  localAddr := conn.LocalAddr().(*net.UDPAddr)
  return localAddr.IP
}

func recvMulticast(addr *net.UDPAddr) {
  // set this socket to listen for multicasts on the specified address
  conn, err := net.ListenMulticastUDP("udp4", nil, addr)
  if err != nil {
    log.Fatalf("Could not join multicast group: %s\n", err)
  }
  conn.SetReadBuffer(maxlen)

  // determine what our IP is
  local := getLocalIP()

  // receive all messages
  for {
    buf := make([]byte, maxlen)
    _, src, err := conn.ReadFromUDP(buf)
    if err != nil {
      log.Printf("Failed to read message: %s\n", err)
    }
    // show the message if it did not come from us
    if !src.IP.Equal(local) {
      fmt.Printf("(%s) %s\n> ", src, string(buf))
    }
  }
}

func sendMulticast(addr *net.UDPAddr) {
  // set this socket to send on the multicast address
  conn, err := net.DialUDP("udp4", nil, addr)
  if err != nil {
    log.Fatalf("Could not set up socket: %s\n", conn)
  }

  // listen for console input
  scanner := bufio.NewScanner(os.Stdin)
  fmt.Print("> ")
  for scanner.Scan() {
    // build message and send it
    str := fmt.Sprintf("%s: %s", *name, scanner.Text())
    conn.Write([]byte(str))
    fmt.Print("> ")
  }
}

func main() {
  // get the multicast address
  addr, err := net.ResolveUDPAddr("udp4", saddr)
  if err != nil {
    log.Fatalf("Could not resolve addr: %s\n", err)
  }

  go recvMulticast(addr)
  sendMulticast(addr)
}
