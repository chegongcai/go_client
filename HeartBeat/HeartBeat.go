package HeartBeat

import (
	"fmt"
	"go_client/BDYString"
	"net"
	"time"
)

func GetMessage(bytes []byte, message chan byte) {
	for _, v := range bytes {
		message <- v
	}
	close(message)
}

func HeartBeat(conn net.Conn, message chan byte, timeout int) {
	select {
	case <-message:
		fmt.Println("message time: ", BDYString.GetTimeStamp())
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	case <-time.After(time.Second * 5):
		fmt.Println(conn.RemoteAddr().String(), "time out")
		conn.Close()
	}
}
