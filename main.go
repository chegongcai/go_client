package main

import (
	"go_client/ClientAndServer"
	"go_client/DeviceAndServer"
)

func main() {
	//go as client
	ClientAndServer.ClientConnetToServer()

	//go as server
	DeviceAndServer.ListenFromDevice()
}
