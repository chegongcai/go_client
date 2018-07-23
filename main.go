package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	addr := "182.254.185.142:8080"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server IP：", conn.RemoteAddr().String())

	fmt.Printf("client IP：%v\n", conn.LocalAddr())

	n, err := conn.Write([]byte("Hello Server"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("send lenght:", n)

	buf := make([]byte, 1024) //定义一个切片的长度是1024。

	n, err = conn.Read(buf) //send lenght

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	fmt.Println(string(buf[:n]))
	conn.Close()
}
