import argparse
import socket
import struct
import threading

# constants
addr = '239.0.0.0'
maxlen = 8192

# variables
name = ''
port = 1024
ttl = struct.pack('b', 1)

# set up flags
parser = argparse.ArgumentParser()
parser.add_argument('-name', help='device name', default="")
parser.add_argument('-port', help='port number', default=1024, type=int)

def getLocalIP():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.connect(('1.1.1.1', 80))
    ip = s.getsockname()[0]
    s.close()
    return ip

def recvMulticast(s):
    s.bind(('', port))
    mreq = struct.pack('4sL', socket.inet_aton(addr), socket.INADDR_ANY)
    s.setsockopt(socket.IPPROTO_IP, socket.IP_ADD_MEMBERSHIP, mreq)
    local = getLocalIP()
    while True:
        buf, src = s.recvfrom(maxlen)
        if src[0] != local:
            print("({}): {}\n> ".format(src[0], buf.decode()), end='')

def sendMulticast(s):
    while True:
        print("> ", end='')
        msg = "{}: {}".format(name, input())
        s.sendto(msg.encode('ascii'), (addr, port))

if __name__ == '__main__':
    # handle flags
    args = parser.parse_args()
    name = args.name
    port = args.port

    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    recvThread = threading.Thread(target=recvMulticast, args=(s,))
    recvThread.start()
    sendMulticast(s)
