package ClientAndServer

import (
	"fmt"
	"go_client/BDYString"
	"go_client/session"
	"net"
	"os"
	"strings"
)

var device_conn, client_conn net.Conn

func ClientConnetToServer() {
	server := "182.254.185.142:8080"
	server_addr, err := net.ResolveTCPAddr("tcp4", server)
	checkErr(err)
	client_conn, err = net.DialTCP("tcp", nil, server_addr)
	checkErr(err)
	go ClientAndServerConn(client_conn)
}

func GetClientConn() net.Conn {
	return client_conn
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func ClientAndServerConn(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("waiting server back msg error: ", err)
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

func ParseServerProtocol(rev_buf string, conn net.Conn) {
	var err error
	var arr_buf, data_buf []string

	arr_buf = strings.Split(rev_buf, "#")             //先分割#
	data_buf = strings.Split(string(arr_buf[5]), ",") //分割;

	fmt.Println(data_buf[0])

	switch data_buf[0] {
	case "ACK^LOCA":
		fmt.Println("get data from go server and then send to device")
		device_conn, err = session.GetConnByID(string(arr_buf[0]))
		if err == nil {
			fmt.Println("device ip: ", device_conn.RemoteAddr().String())
			_, err = device_conn.Write([]byte(rev_buf))
		}
		break
	}
	fmt.Println("****************************************************************************************")
	if err != nil {
		return
	}
	fmt.Println("****************************************************************************************")
}
