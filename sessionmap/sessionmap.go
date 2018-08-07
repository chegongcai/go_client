package sessionmap

import (
	"fmt"
	"net"
)

type Session struct {
	id   string
	conn net.Conn
}

var sessionmap = make(map[string]Session)

func AddNewSessionMap(id string, conn net.Conn) {
	sessionmap[id] = Session{id, conn}
}

func DeleteOneSessionMap(id string) {
	delete(sessionmap, id)
}

func GetConnByIDMap(id string) (net.Conn, bool) {
	var conn net.Conn
	session, ok := sessionmap[id]
	if ok {
		fmt.Println("get conn", session.conn)
		return session.conn, true
	} else {
		return conn, false
	}
}

func CheckSessionMap(id string) bool {
	_, ok := sessionmap[id]
	if ok {
		return true
	} else {
		return false
	}
}
