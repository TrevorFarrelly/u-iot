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

  Having a framework is great and all, but what use is a smart home network without smart home _devices_? u-iot will provide build tools for easily installing a lightweight Raspbian distribution based on [u-root](https://github.com/u-root/u-root). With the use of a GPIO library and some basic circuitry skills, one can build a fully functional smart light, outlet, fan, plant waterer, fish feeder, or whatever else they can dream of.

  u-root has its own set of build tools, wrapping them into an easy-to-use u-iot installer will be the main goal of this portion of the project. u-root is also very barebones, so I imagine lots of experimentation will be necessary to make everything function out-of-the-box.

### Deliverables
1. Communication - March 6
  * send and receive multicasts between multiple devices.
  * implement bootstrapping process via protobufs and gRPC.
  * extend protobuf to include function calls and parameter passing.
  * write detailed communication protocol documentation, finalize RPCs.

  __Result__: toy program that can run on multiple devices, demonstrating communication between them. Documentation of the protocol for any future developers.

2. Library - March 20
  * generate Go protobuf code and build any necessary wrappers (parameter types, etc).
  * polish bootstrapping process used in test programs, add any necessary interfaces.
  * Implement a second library in another popular language (TBD, most likely Python).

  __Result__: Go and Python libraries that can be imported and used in any program, with functional language-agnostic communication.

3. Raspberry Pi - End of Semester
  * Boot manually installed u-root installation, write install script to wrap process.
  * Incorporate u-iot libraries and dependencies into installation script.
  * Experiment with internet - ethernet, WiFi, etc.
  * Integrate internet configuration into install.
  * Build basic example devices and apps:
    * CLI interface
    * Repurpose old, broken fan
    * RGB LED strip (?)

  __Result__: Three devices that communicate with each other and can be controlled by each other.

### More Detailed Ideas
This section is mostly a place for me to get minute implementation details down on paper before I forget them.

* Basic API (Go) - I imagine most other languages with actual OOP features will have a slightly different design to take advantage of it.

u-iot's interface should be built with ease in mind. I want someone to be able to build a basic Go program with a some functions, import `uiot-go`, add a couple lines, and have a working u-iot device.
  * `uiot.Func` - This type associates a Go `func()` with a name and a list of parameters. This is how one defines which functions can be called by other devices on the network.
  * `uiot.Param` - This type provides a "universal" parameter that can be used for most functions and devices. It will most likely represent a range of possible integers, with aliases for potentially common types:
    * `uiot.BoolParam` = `uiot.Param(0,1)`
    * `uiot.LightParam` = `uiot.Param(0,255)`
    * `uiot.RGBParam` = `[ uiot.LightParam, uiot.LightParam, uiot.LightParam ]`
  * `uiot.Register(name string, funcs []uiot.Function)` - this call tells u-iot what name to use and all of the functions it has. Will need to be called before bootstrapping.
  * `uiot.Bootstrap(retry string)` - Attempt to connect to the network and build a "database" of known devices asynchronously. Takes a [`time.Duration`](https://pkg.go.dev/time?tab=doc#Duration)-parseable parameter for how often to refresh the database. May return a way to access that database for UI programs, not sure exactly what that would be yet.
