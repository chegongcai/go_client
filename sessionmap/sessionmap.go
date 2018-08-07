package sessionmap

import (
	"net"
)

type SessionMap struct {
	id   string
	conn net.Conn
}

var sessionmap = make(map[string]SessionMap)

func AddNewSessionMap(id string, conn net.Conn) {
	sessionmap[id] = SessionMap{id, conn}
}

func DeleteOneSessionMap(id string) {
	delete(sessionmap, id)
}

func GetConnByIDMap(id string) (net.Conn, bool) {
	//var conn net.Conn
	session, ok := sessionmap[id]
	if ok {
		return session.conn, true
	} else {
		return nil, false
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
