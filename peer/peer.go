package peer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
	"vayeate-api/logger"
)

// timeout connections after 60s of inactivity
const timeout = 60 * time.Second

var UnauthorizedError = errors.New("peer: unauthorized")

type Peer struct {
	ID        string    `json:"id"`
	StartedAt time.Time `json:"started_at"`

	closed    bool
	opened    bool
	log       *logger.Logger
	pingTimer *time.Timer
	reader    *bufio.Reader
	conn      net.Conn
}

func FromConnection(nodeId string, username string, password string, conn net.Conn) (*Peer, error) {
	reader := bufio.NewReader(conn)
	now := time.Now()
	self := Peer{"", now, false, false, nil, nil, reader, conn}
	self.pingTimer = time.AfterFunc(timeout, onTimeout(&self))
	m, err := self.Read()

	if err != nil {
		return nil, err
	}

	self.ID = m.Headers.FromID
	self.log = logger.New(fmt.Sprintf("vayeate:peer:%s", self.ID))
	err = self.Write(NewOpenSuccessMessage(nodeId, username, password))
	return &self, err
}

func Connect(nodeId string, username string, password string, address string) (*Peer, error) {
	conn, err := net.Dial("tcp", address)

	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)
	now := time.Now() // to-do: get started_at from peer
	self := Peer{"", now, false, false, nil, nil, reader, conn}
	self.pingTimer = time.AfterFunc(timeout, onTimeout(&self))
	err = self.Write(NewOpenMessage(nodeId, username, password))

	if err != nil {
		return nil, err
	}

	m, err := self.Read()

	if err != nil {
		return nil, err
	}

	self.ID = m.Headers.FromID
	self.log = logger.New(fmt.Sprintf("vayeate:peer:%s", self.ID))
	return &self, nil
}

func (self *Peer) Close() {
	self.closed = true
	self.opened = false
	self.conn.Close()
	self.pingTimer.Stop()
}

func (self *Peer) Closed() bool {
	return self.closed
}

func (self *Peer) Opened() bool {
	return self.opened
}

func (self *Peer) GetRemoteAddress() string {
	return self.conn.RemoteAddr().String()
}

func (self *Peer) Read() (*Message, error) {
	m, err := DecodeMessage(self.reader)

	if err == io.EOF {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if !self.opened && !m.IsOpen() && !m.IsOpenSuccess() {
		return nil, UnauthorizedError
	}

	self.opened = true
	self.pingTimer.Reset(timeout)
	return m, nil
}

func (self *Peer) Write(m *Message) error {
	b, err := m.Encode()

	if err != nil {
		return err
	}

	_, err = self.conn.Write(b)
	return err
}

func onTimeout(self *Peer) func() {
	return func() {
		self.Close()
		self.log.Infoln("closed due to inactivity")
	}
}
