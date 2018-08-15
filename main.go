package main

import (
	"go_client/ClientAndServer"
	"go_client/DeviceAndServer"
	//"go_client/gomysql"
)

func main() {
	//Init mysql service
	//gomysql.Init()

	//go as client
	ClientAndServer.ClientConnetToServer()

	//go as server
	DeviceAndServer.ListenFromDevice()
}
