# Protocols

u-iot communicates using a number of technologies. This README documents each component and the steps required for bootstrapping and sending commands, both for documentation and implementation in other languages in the future.

### Protocol Buffers and gRPC

u-iot's main functionality is in `uiot.proto` in this directory. Described there are the RPC interfaces and messages that are sent between devices.

#### Bootstrapping RPCs

* `DevInfo` is sent and received on startup, containing a device's IP, port, human-readable meta information, and a set of `FuncDef`s describing the functions it performs.

* `FuncDef` describes a function a device performs. Contains an ID, human-readable name, and a list of `ParamDef`s, each representing a minimum and maximum value the function expects.

#### Triggering RPCs

* TBD

### Multicasting

gRPC communication requires knowing a remote device's IP and port number. To get that information, we send UDP multicast packets to a group that all u-iot devices are listening on. Each packet contains the port the device is using for RPC communication, encoded in big-endian.

## Bootstrapping Process

When a u-iot program starts, it bootstraps itself with all other u-iot programs on the network. The steps a new program takes are as follows:

```
local device                           router                    remote device(s)
     |                                   |                              |
     |     send RPC port over UDP to     |  >---->---->---->---->---->  |
  1. |  >---->---->---->---->---->---->  |    forward to all devices    |
     |   multicast addr 239.0.0.0:1024   |  >---->---->---->---->---->  |
     |                                   |                              |
     |                                                                  |
     |         send Bootstrap RPC with remote device information        |
  2. |  <----<----<----<----<----<----<-----<----<----<----<----<----<  |
     |                                                                  |
     |        respond to Bootstrap with local device information        |
  3. |  >---->---->---->---->---->---->----->---->---->---->---->---->  |
     |                                                                  |
```
1. All u-iot devices have a UDP server listening on `239.0.0.0:1024`. When a new device starts up, it sends the port its gRPC server is using to this address, encoded in 2 network-order (big-endian) bytes. Since this IP is a multicast address, the router copies the sent packet for every other device on the network.

2. Using the gRPC port and source IP received in step 1, all remote devices send a Bootstrap RPC with their device information to the new device. The new device saves this info to their set of known devices.

3. The new device responds to the remote RPCs with its own device info, saving it to their sets. All existing devices know about the new device, and the new device knows about all existing devices.
