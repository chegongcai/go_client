package session

import (
	"fmt"
	"net"
	"strings"
	"syscall"
)

type Session struct {
	id   string
	conn net.Conn
}

type SessionP struct {
	session []*Session
}

var bc = &SessionP{}

func AddNewSession(id string, conn net.Conn) {
	session := &Session{id, conn}
	bc.session = append(bc.session, session)
}

func GetConnByID(id string) (net.Conn, error) {
	var conn net.Conn
	for _, block := range bc.session {
		fmt.Println(block.id)
		if strings.Contains(id, block.id) {
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
