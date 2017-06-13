package telnet

import (
	"bufio"
	"bytes"
	"connutil"
	"fmt"
	"net"
	"time"
)

const (
	newline        = "\n"
	loginPrompt    = "login: "
	passwordPrompt = "password: "
	username       = "itsthebishop"
	password       = ""
	systemPrompt   = "NET> "
	telnetPort     = "23"

	ip = "10.4.3.240"

	dialTimeout = 1 * time.Second
)

var (
	monitoringCommand = []byte("#MONITORING,25,ON")

	messageDelimiter = []byte("\r\n")

	badLogin = []byte("bad login\r\n")
)

type Conn struct {
	net.Conn
	reader    *bufio.Reader
	connected bool
}

func NewConnection() (*Conn, error) {
	conn, err := net.DialTimeout("tcp", ip+":"+telnetPort, dialTimeout)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(conn)

	err = setupSession(conn, username, password, reader)
	if err != nil {
		return nil, err
	}

	connWrapper := &Conn{
		Conn:      conn,
		reader:    reader,
		connected: true,
	}

	go connWrapper.readOnConnection()

	return connWrapper, nil

}

func setupSession(conn net.Conn, username string, password string, reader *bufio.Reader) error {
	if err := connutil.SkipUntil(reader, loginPrompt); err != nil {
		return err
	}

	if _, err := conn.Write([]byte(username)); err != nil {
		return err
	}

	if _, err := conn.Write(messageDelimiter); err != nil {
		return err
	}

	if err := connutil.SkipUntil(reader, passwordPrompt); err != nil {
		return err
	}

	if _, err := conn.Write([]byte(password)); err != nil {
		return err
	}

	if _, err := conn.Write(messageDelimiter); err != nil {
		return err
	}

	response, err := connutil.ReadUntil(reader, systemPrompt, newline)
	if err != nil {
		return err
	}

	if bytes.Contains(response, badLogin) {
		return fmt.Errorf("Error badlogin")
	}

	fmt.Printf("Login successful for remote client %v.\n", conn.RemoteAddr())

	if _, err := conn.Write(monitoringCommand); err != nil {
		return err
	}

	if _, err := conn.Write(messageDelimiter); err != nil {
		return err
	}

	return nil
}

func (conn *Conn) readOnConnection() {
	for {

		line, err := conn.reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error while reading on connection, : %v\n", err)
			conn.Close()
		}
		fmt.Printf("Line Read: %s\n", line)
	}
}

func (conn *Conn) Connected() bool {
	return conn.connected
}
