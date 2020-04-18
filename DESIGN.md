# u-iot

u-iot is a protocol and framework for building custom smart home applications. Intended to provide complete freedom in device functionality, u-iot will support applications on host operating systems, as well as embedded devices, abilities, and parameters on Raspberry Pi.

u-iot will be primarily written in Go, but the protocol will be designed to support most modern languages. See Libraries below for more details.

### Structure

u-iot has three main components: communication protocols, libraries, and Pi installers.

* __Communication Protocols__

  u-iot's primary communication stack will be built on [gRPC](https://grpc.io/). gRPC enables efficient, decentralized, peer-to-peer communication, and it is language-agnostic, allowing for u-iot libraries to be built for most modern languages. Generic devices, abilities, and parameters will be defined in this layer.

  RPCs are great for quickly sending messages, but connecting to a u-iot RPC network will require bootstrapping with a known IP address. In an effort to avoid this, device discovery will be done with [UDP multicasting](https://en.wikipedia.org/wiki/Multicast#IP_multicast). When a new device first powers on, it will send a multicast request for other u-iot devices and their IP addresses. From there the new device can request other information from every other device via RPCs.

* __Libraries__

  u-iot libraries will be language specific, facilitating communication and providing an interface for users to build their own apps and devices.

  The first library will be built for Go, but all [gRPC-supported languages](https://grpc.io/docs/) could have libraries built for them. I predict at least one other would be built over the course of development.

* __Raspberry Pi Integration__

  Having a framework is great and all, but what use is a smart home network without smart home _devices_? u-iot will provide build tools for easily installing a lightweight Raspbian distribution based on [u-root](https://github.com/u-root/u-root). With the use of a GPIO library (such as [go-rpio](https://github.com/stianeikeland/go-rpio))and some basic circuitry skills, one can build a fully functional smart light, outlet, fan, plant waterer, fish feeder, or whatever else they can dream of.

  u-root has its own set of build tools, wrapping them into an easy-to-use u-iot installer will be the main goal of this portion of the project. u-root is also very barebones, so I imagine lots of experimentation will be necessary to make everything function out-of-the-box.

### Deliverables

_Februrary_
  * Send and receive multicasts between numerous devices across the network.
  * Define bootstrapping process in protobufs and gRPC.
    * Extend protobuf to include function calls and parameter passing.
  * Implement toy program that demonstrates the entire bootstrapping process.
  * Write documentation on communication protocol to ease implementation in other languages.

_March_
  * Build uiot-go library.
    * Extract bootstrapping process from toy program.
    * Expose function signature definition interface as an easy-to-use API.
    * Reimplement toy program using the new library.
  * Build uiot-python to demonstrate interoperability between languages.
  * Experiment with a repeatable process for installing u-root on Raspberry Pi.
    * I have a B+, Zero W, and 4 to test on.

_April_
  * Wrap process in installer script.
  * Integrate internet configuration into the installer.
  * Allow for user to specify files to include and which programs to run on startup.
  * Build basic example devices.
    * CLI interface
    * LED strip
    * Desk fan
