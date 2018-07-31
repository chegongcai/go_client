package main

import (
	"fmt"
	"go_client/BDYString"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var device_conn net.Conn

//182.254.185.142  8080
const version = 0 // 0 for debug
var SerialNum int
var send_test int = 0

//for connect conn
type Session struct {
	id   string
	conn net.Conn
}

type SessionP struct {
	session []*Session
}

var bc = &SessionP{}

func NewSession(id string, conn net.Conn) *Session {
	session := &Session{id, conn}
	return session
}

func (bc *SessionP) AddSession(id string, conn net.Conn) {
	newSession := NewSession(id, conn)
	bc.session = append(bc.session, newSession)
}

func GetConnByID(id string) (net.Conn, error) {
	var conn net.Conn
	for _, block := range bc.session {
		fmt.Println(block.id)
		fmt.Println(block.conn.RemoteAddr().String())
		if strings.Contains(id, block.id) {
			fmt.Println("get conn")
			return block.conn, nil
		}
	}
	return conn, syscall.EINVAL
}

func DeleteConnByID(id string) {
	for index, block := range bc.session {
		if strings.Contains(id, string(block.id)) {
			bc.session = append(bc.session[:index], bc.session[index+1:]...)
		}
	}
}

func CheckSession(id string) error {
	for _, block := range bc.session {
		if strings.Contains(id, string(block.id)) {
			return nil
		}
	}
	return syscall.EINVAL
}

//for connect conn

func main() {
	//server
	service := ":8080"
	//testbuf()
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

	//client
	ClientConnetToServer()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go DeviceAndServerConn(conn)
	}
}

func ClientConnetToServer() net.Conn {
	server := "182.254.185.142:8080"
	server_addr, err := net.ResolveTCPAddr("tcp4", server)
	checkErr(err)
	client_conn, err := net.DialTCP("tcp", nil, server_addr)
	checkErr(err)
	go ClientAndServerConn(client_conn)
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
		fmt.Println("time: ", GetTimeStamp())
		fmt.Println("receive data from server: ", string(buffer[:n]))
		if buffer[n-1] != '$' {
			return
		}
		rev_buf := string(buffer[0 : n-1]) //delete the tail #
		ParseServerProtocol(rev_buf, conn) //do protocol parse
	}
}

func DeviceAndServerConn(conn net.Conn) {
	defer conn.Close()

	var buf [1024]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println("conn close", n, conn.RemoteAddr().String())
			DeleteConnByID(conn.RemoteAddr().String()) //释放断开的链接绑定
			return
		}
		rAddr := conn.RemoteAddr()
		fmt.Println("****************************************************************************************")
		fmt.Println("device ip: ", rAddr.String())
		fmt.Println("time: ", GetTimeStamp())
		fmt.Println("rev data: ", string(buf[0:n]))
		if buf[n-1] != '$' {
			return
		}
		rev_buf := string(buf[0 : n-1])    //delete the tail #
		ParseDeviceProtocol(rev_buf, conn) //do protocol parse
	}
}

func GetZone() string {
	local, _ := time.LoadLocation("Local")
	local_str := fmt.Sprintf("%s", time.Now().In(local))
	buf := []byte(local_str)
	return string(buf[32:33])
}

func GetTimeStamp() string {
	buf := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	return buf
}

func GetTimeStampForSYNC() string {
	zone, _ := strconv.Atoi(GetZone())
	buf := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour()-zone, time.Now().Minute(), time.Now().Second())
	return buf
}

func testbuf() {

}

func ParseDeviceProtocol(rev_buf string, conn net.Conn) {
	var err error
	var arr_buf, data_buf, comand_buf []string

	arr_buf = strings.Split(rev_buf, "#")                //先分割#
	data_buf = strings.Split(string(arr_buf[4]), ";")    //分割;
	comand_buf = strings.Split(string(data_buf[0]), ":") //分割;

	fmt.Println(comand_buf[0])
	serial_num := string(arr_buf[2])
	imei := string(arr_buf[1])

	SerialNum = BDYString.HexString2Int(serial_num)

	err_check := CheckSession(conn.RemoteAddr().String()) //检查是否是已经绑定过的
	if err_check != nil {
		bc.AddSession(conn.RemoteAddr().String(), conn) //绑定ip与conn
	}

	switch comand_buf[0] {
	case "LOCA":
		//parse data
		switch comand_buf[1] {
		case "W":
			alert := BDYString.GetBetweenStr(rev_buf, "ALERT", ";")
			status := BDYString.GetBetweenStr(rev_buf, "STATUS", ";")
			//wifi := BDYString.GetBetweenStr(rev_buf, "WIFI", "$")
			fmt.Println(status)
			fmt.Println(alert)
			//fmt.Println(wifi)
			break
		case "G":
			GPS_DATA := BDYString.GetBetweenStr(rev_buf, "GDATA", ";")
			alert := BDYString.GetBetweenStr(rev_buf, "ALERT", ";")
			status := BDYString.GetBetweenStr(rev_buf, "STATUS", ";")
			fmt.Println(status)
			fmt.Println(alert)
			fmt.Println(GPS_DATA)
			break
		case "L":
			LBS_DATA := BDYString.GetBetweenStr(rev_buf, "CELL", ";")
			alert := BDYString.GetBetweenStr(rev_buf, "ALERT", ";")
			status := BDYString.GetBetweenStr(rev_buf, "STATUS", ";")
			fmt.Println(status)
			fmt.Println(alert)
			fmt.Println(LBS_DATA)
			break
		}
		fmt.Println("send data to server")
		buf := fmt.Sprintf("%s#S168#%s#%s#0009#ACK^LOCA,$", conn.RemoteAddr().String(), imei, serial_num)
		_, err = ClientConnetToServer().Write([]byte(buf)) //send to server
		break
	case "B2G":
		//parse data
		var lbs_buf []string
		var lbs_int [4]int
		lbs_buf = strings.Split(string(comand_buf[1]), ",") //分割;
		for i := 0; i < 4; i++ {
			lbs_int[i] = BDYString.HexString2Int(string(lbs_buf[i]))
		}
		fmt.Println(lbs_int)
		//printf data

		//send data  //22.529793,113.952744
		buf := fmt.Sprintf("S168#%s#%s#0028#ACK^B2G,22.529793,113.952744$", imei, serial_num)
		fmt.Println("device ip: ", conn.RemoteAddr().String())
		fmt.Println("send data: ", buf)
		_, err = conn.Write([]byte(buf))
		break

	case "SYNC":
		////parse data
		var buf string
		if comand_buf[1] == "0000" { //收到登录包
			//bc.AddSession(conn.RemoteAddr().String(), conn) //绑定ip与conn
			buf = fmt.Sprintf("S168#%s#%s#0023#ACK^SYNC,%s$", imei, serial_num, GetTimeStampForSYNC())
		} else {
			buf = fmt.Sprintf("S168#%s#%s#0009#ACK^SYNC,$", imei, serial_num)
		}
		fmt.Println("send data: ", buf)
		_, err = conn.Write([]byte(buf))
		break
	}
	if err != nil {
		return
	}

	if send_test == 1 {
		send_test = 1
		SerialNum++
		buf := fmt.Sprintf("S168#%s#%s#0009#GSENSOR,1$", imei, BDYString.Int2HexString(SerialNum))
		fmt.Println("send data: ", buf)
		_, err = conn.Write([]byte(buf))
	}
	fmt.Println("****************************************************************************************")
	//device_conn = conn
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
		device_conn, err = GetConnByID(string(arr_buf[0]))
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

/*
package main

import (
	"fmt"
	"strings"
	"syscall"
)

type Session struct {
	Data []byte
	ID   []byte
}

type SessionP struct {
	session []*Session
}

var bc = &SessionP{}

func NewSession(data string, id []byte) *Session {
	session := &Session{[]byte(data), id}
	return session
}

func (bc *SessionP) AddSession(data string, id string) {
	newSession := NewSession(data, []byte(id))
	bc.session = append(bc.session, newSession)
}

func GetIDByData(data string) []byte {
	for _, block := range bc.session {
		if strings.Contains(data, string(block.Data)) {
			return block.ID
		}
	}
	return nil
}

func DeleteConnByID(data string) {
	fmt.Println(bc.session)
	for index, block := range bc.session {
		fmt.Println(index)
		if strings.Contains(data, string(block.Data)) {
			bc.session = append(bc.session[:index], bc.session[index+1:]...)
			fmt.Println(bc.session)
		}
	}
}

func CheckSession(data string) error {
	for index, block := range bc.session {
		fmt.Println(index)
		if strings.Contains(data, string(block.Data)) {
			fmt.Println("the same data")
			return nil
		}
	}
	return syscall.EINVAL
}

func InitBlock() {
	bc.AddSession("Send 1", "123")
	err := CheckSession("Send 1")
	if err != nil {
		bc.AddSession("Send 1", "123")
	}
	bc.AddSession("Send 2", "456")
	bc.AddSession("Send 3", "789")
}

func main() {
	InitBlock()
	fmt.Println(bc.session)
	DeleteConnByID("Send 2")
	fmt.Println(string(GetIDByData("Send 2")))
}
*/
