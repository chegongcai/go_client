package main

import (
	"go_client/ClientAndServer"
	"go_client/DeviceAndServer"
	"go_client/gomysql"
)

func main() {
	gomysql.Init()
	//go as client
	ClientAndServer.ClientConnetToServer()

	//go as server
	DeviceAndServer.ListenFromDevice()
}
