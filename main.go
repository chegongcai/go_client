package main

import (
	"fmt"
	"go_client/ClientAndServer"
	"go_client/DeviceAndServer"
	"net"
	"os"
	"time"
)

//182.254.185.142  8080
const version = 0 // 0 for debug
var SerialNum int
var send_test int = 0

func main() {
	//server
	service := ":8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

	//client
	ClientAndServer.ClientConnetToServer()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		conn.SetDeadline(time.Now().Add(time.Duration(60) * time.Second))
		go DeviceAndServer.DeviceAndServerConn(conn)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
