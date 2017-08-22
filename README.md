
# p2p-port-forward

[![Build Status](https://travis-ci.org/selvakn/p2p-port-forward.svg?branch=master)](https://travis-ci.org/selvakn/p2p-port-forward) [![Go Report Card](https://goreportcard.com/badge/github.com/selvakn/p2p-port-forward)](https://goreportcard.com/report/github.com/selvakn/p2p-port-forward)

Command line utility to forward ports between two hosts across different networks/subnets, in a peer-to-peer fashion using [zerotier](zerotier.com), _without_ sudo privileges

### Installation

On both the hosts.

1) Homebrew

    `brew install https://raw.githubusercontent.com/selvakn/homebrew-core/master/Formula/p2p-port-forward.rb`

    Or 

2) Using `go get`

    `go get https://github.com/selvakn/p2p-port-forward`

    Or
    
3) Single binaries
    * Download the single binary for your platform from [releases](https://github.com/selvakn/p2p-port-forward/releases)
    * `mv p2p-port-forward-<platform/arch> p2p-port-forward`
    * `chmod +x p2p-port-forward`

### Usage

```
usage: p2p-port-forward [<flags>]

Flags:
      --help                   Show context-sensitive help (also try --help-long and --help-man).
  -n, --network="8056c2e21c000001"
                               zerotier network id
  -f, --forward-port="22"      port to forward (in listen mode)
  -a, --accept-port="2222"     port to accept (in connect mode)
  -u, --use-udp                UDP instead of TCP (TCP default)
  -c, --connect-to=CONNECT-TO  server (zerotier) ip to connect
      --version                Show application version.

```

On host1 (where the service or any port should be exposed)

`p2p-port-forward [-n 8056c2e21c000001] [-f 22]`
 
On host2 (where we want to access the port)

`p2p-port-forward -c <server-ip-from-output-of-server> [-a 2222]`


Then connect to port 2222 on host2, to reach host1:22. Ex:
 
 `ssh -p 2222 user@localhost`
 