package ClientAndServer

import (
	"fmt"
	"go_client/BDYString"
	"go_client/sessionmap"
	"net"
	"os"
	"strings"
	"time"
)

var device_conn, client_conn net.Conn
var connect_status int = 0

/*
description:
go as a client
start connect to the server(JAVA server)
*/
func ClientConnetToServer() {
	server := "182.254.185.142:8080"
	server_addr, err := net.ResolveTCPAddr("tcp4", server)
	checkErr(err)
	client_conn, err = net.DialTCP("tcp", nil, server_addr)
	if err != nil {
		fmt.Println("Reconnect to server...")
		connect_status = 0
		time.AfterFunc(5*time.Second, ClientConnetToServer)
	} else {
		fmt.Println("Connect to server OK!")
		connect_status = 1
		go ClientAndServerConn(client_conn)
	}
}

func GetConnectStatus() int {
	return connect_status
}

/*
description:
return the network connection
*/
func GetClientConn() net.Conn {
	return client_conn
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

/*
description: get the network connection data, and other messege
input: net.conn
return: no
*/
func ClientAndServerConn(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("waiting server back msg error: ", err)
			ClientConnetToServer()
			return
		}
		fmt.Println("****************************************************************************************")
		fmt.Println("server ip: ", conn.RemoteAddr().String())
		fmt.Println("time: ", BDYString.GetTimeStamp())
		fmt.Println("receive data from server: ", string(buffer[:n]))
		if buffer[n-1] != '$' {
			return
		}
		rev_buf := string(buffer[0 : n-1]) //delete the tail #
		ParseServerProtocol(rev_buf, conn) //do protocol parse
	}
}

/*
description: do the protocol parse, and then send the data to the device
input: rev_buf->data from conn.read
	   conn->the network connection
return: if error, will do the return
*/
func ParseServerProtocol(rev_buf string, conn net.Conn) {
	var err error
	var arr_buf, data_buf []string
	var ok bool

	arr_buf = strings.Split(rev_buf, "#")             //parse "#"
	data_buf = strings.Split(string(arr_buf[5]), ",") //parse ";"

	fmt.Println(data_buf[0])

	switch data_buf[0] {
	case "ACK^LOCA":
		fmt.Println("get data from go server and then send to device")
		device_conn, ok = sessionmap.GetConnByIDMap(string(arr_buf[0]))
		if ok == true {
			fmt.Println("device ip: ", device_conn.RemoteAddr().String())
			_, err = device_conn.Write([]byte(rev_buf))
		}
		break
	}
	if err != nil {
		return
	}
	fmt.Println("****************************************************************************************")
}
