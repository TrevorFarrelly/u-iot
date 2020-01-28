# u-iot

u-iot is a protocol and framework for building custom smart home applications. Intended to provide complete freedom in device functionality, u-iot will support applications on host operating systems, as well as embedded devices, abilities, and parameters on Raspberry Pi.

u-iot will be written in Go.

### Structure

u-iot has three main components: communication protocols, libraries, and Pi installers.

* __Communication Protocols__

  u-iot's primary communication stack will be built on [gRPC](https://grpc.io/). gRPC enables efficient, decentralized, peer-to-peer communication, and it is language-agnostic, allowing for u-iot libraries to be built for most modern languages. Generic devices, abilities, and parameters will be defined in this layer.

  RPCs are great for quickly sending messages, but connecting to a u-iot RPC network will require bootstrapping with a known IP address. In an effort to avoid this, device discovery will be done with [UDP multicasting](https://en.wikipedia.org/wiki/Multicast#IP_multicast). When a new device first powers on, it will send a multicast request for other u-iot devices and their IP addresses. From there the new device can request other information from every other device via RPCs.

* __Libraries__

  u-iot libraries will be language specific, facilitating communication and providing an interface for users to build their own apps and devices.

  The first library will be built for Go, but all [gRPC-supported languages](https://grpc.io/docs/) could have libraries built for them. I predict at least one other would be built over the course of development.

* __Raspberry Pi Integration__

  Having a framework is great and all, but what use is a smart home network without smart home _devices_? u-iot will provide build tools for easily installing a lightweight Raspbian distribution based on [u-root](https://github.com/u-root/u-root). With the use of a GPIO library and some basic circuitry skills, one can build a fully functional smart light, outlet, plant waterer, fish feeder, or whatever else they can dream of.

  u-root has its own set of build tools, wrapping them into an easy-to-use u-iot installer will be the main goal of this portion of the project. It is also very barebones, so I imagine lots of experimentation will be necessary (and thus most development time will be spent on this portion of the project).

### Deliverables
1. Communication - March 6
  * send and receive multicasts between multiple devices.
  * implement bootstrapping process via protobufs and gRPC.
  * extend protobuf to include function calls and parameter passing.
  * write detailed communication protocol documentation, finalize RPCs.
2. Library - March 27
  * generate Go protobuf code and build any necessary wrappers (parameter types, etc).
  * polish bootstrapping process used in test programs, add any necessary interfaces.
  * Implement a second library in another popular language (Java? Python? TBA).
3. Raspberry Pi - End of Semester
  * Boot manually installed u-root installation, write install script to wrap process.
  * Incorporate u-iot libraries and dependencies into installation script.
  * Experiment with internet - ethernet, WiFi, etc.
  * Integrate internet configuration into install.
  * Build basic example devices and apps:
    * CLI interface
    * Repurpose old, broken fan
    * RGB LED strip (?)

### Notes and Ramblings
* Possible Ability Parameters - Since protocol buffers do not have traditional generics/inheritance, we need a single message type that encapsulates most necessary parameters as simply as possible. Right now, two integers seems like the best bet:
  * momentary - no parameters   - none
  * toggle    - on/off          - 0 to 1
  * slider    - min to max      - (min) to (max)
  * RGB       - 3 0-255 sliders - 0-255, 0-255, 0-255
