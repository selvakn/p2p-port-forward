package main

import (
	"net"
	"syscall"
	"./forwarder"
	"./listener"
	"./utils"
)

const PORT = 50718 // 7878

func main() {
	utils.InitZT()

	if len(forwarder.GetOtherIP()) == 0 {
		sockfd := listener.BindAndListen(PORT)

		go func() {
			for {
				newSockfd := listener.Accept(sockfd)
				go listener.HandleIncoming(newSockfd)
			}
		}()

		<-utils.SetupCleanUpOnInterrupt(func() {
			syscall.Close((int)(sockfd))
		})

	} else {
		ln, _ := net.Listen("tcp", ":2222")

		go func() {
			for {
				conn, err := ln.Accept()
				if err == nil {
					go forwarder.HandleOutgoing(conn, PORT)
				}
			}
		}()

		<-utils.SetupCleanUpOnInterrupt(func() {
			ln.Close()
		})

	}

}
