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
	session []Session
}

var bc = SessionP{}

/*
description: add one new session
input: id->ip string
	   conn->the network connection
return: no
*/
func AddNewSession(id string, conn net.Conn) {
	session := Session{id, conn}
	bc.session = append(bc.session, session)
}

/*
description: get the network connection by the ip string
input: id->ip string
return: net.conn->the network connection
	    error->the result of get,nil mean can get the conn
*/
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

/*
description: delete the network connection by the ip string
input: id->ip string
return: no
*/
func DeleteConnByID(id string) {
	for index, block := range bc.session {
		if strings.Contains(id, string(block.id)) {
			bc.session = append(bc.session[:index], bc.session[index+1:]...)
		}
	}
}

/*
description: check the ip string is already binding or not
input: id->ip string
return: nil mean had already binding
*/
func CheckSession(id string) error {
	for _, block := range bc.session {
		if strings.Contains(id, string(block.id)) {
			return nil
		}
	}
	return syscall.EINVAL
}
