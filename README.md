# p2p-port-forward

Command line utility to forward ports between two hosts across different networks/subnets, in a peer-to-peer fashion using [zerotier](zerotier.com).
, _without_ sudo privileges

### Installation

On both the hosts

1) Using `go get`

    `go get https://github.com/selvakn/p2p-port-forward`

    Or
    
2) Single binaries
    * Download the single binary for your platform from [releases](https://github.com/selvakn/p2p-port-forward/releases)
    * `mv p2p-port-forward-<platform/arch> p2p-port-forward`

### Usage

On host1 (where the service or any port should be exposed)

`p2p-port-forward [-n 8056c2e21c000001] [-f 22]`
 
On host2 (where we want to access the port)

`p2p-port-forward -c <server-ip-from-output-of-server> [-a 2222]`


Then connect to port 2222 on host2, to reach host1:22. Ex:
 
 `ssh -p 2222 user@localhost`
 