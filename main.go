package main

import (
	"chetest/BDYString"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//182.254.185.142  8080
const version = 0 // 0 for debug
var SerialNum int
var send_test int = 0

func main() {
	//beego.Run()
	service := ":8080"
	//testbuf()
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [1024]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		rAddr := conn.RemoteAddr()
		fmt.Println("****************************************************************************************")
		fmt.Println("client ip: ", rAddr.String())
		fmt.Println("time: ", GetTimeStamp())
		fmt.Println("rev data: ", string(buf[0:n]))
		if buf[n-1] != '$' {
			return
		}
		rev_buf := string(buf[0 : n-1]) //delete the tail #
		ParseProtocol(rev_buf, conn)    //do protocol parse
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
	buf := "S168#358511029674984#000e#0012#SYNC:0003;CLOSED:1"

	var arr_buf, data_buf, comand_buf []string

	arr_buf = strings.Split(buf, "#")                    //先分割#
	data_buf = strings.Split(string(arr_buf[4]), ";")    //分割;
	comand_buf = strings.Split(string(data_buf[0]), ":") //分割:

	fmt.Println(comand_buf[1])
}

func ParseProtocol(rev_buf string, conn net.Conn) {
	var err error
	var arr_buf, data_buf, comand_buf []string

	//fmt.Println("Receive from client", rev_buf)

	arr_buf = strings.Split(rev_buf, "#")                //先分割#
	data_buf = strings.Split(string(arr_buf[4]), ";")    //分割;
	comand_buf = strings.Split(string(data_buf[0]), ":") //分割;

	fmt.Println(comand_buf[0])
	serial_num := string(arr_buf[2])
	imei := string(arr_buf[1])

	SerialNum = BDYString.HexString2Int(serial_num)

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
		//printf data  //len([]rune(buf))-27
		//send data
		buf := fmt.Sprintf("S168#%s#%s#0009#ACK^LOCA,$", imei, serial_num)
		fmt.Println("send data: ", buf)
		_, err = conn.Write([]byte(buf))
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
		fmt.Println("send data: ", buf)
		_, err = conn.Write([]byte(buf))
		break

	case "SYNC":
		////parse data
		var buf string
		if comand_buf[1] == "0000" {
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
	if send_test == 0 {
		send_test = 1
		SerialNum++
		buf := fmt.Sprintf("S168#%s#%s#0009#GSENSOR,1$", imei, BDYString.Int2HexString(SerialNum))
		fmt.Println("send data: ", buf)
		_, err = conn.Write([]byte(buf))
	}
	fmt.Println("****************************************************************************************")
}
