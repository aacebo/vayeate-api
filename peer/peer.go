package peer

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
	"vayeate-api/logger"

	"github.com/google/uuid"
)

// timeout connections after 60s of inactivity
const timeout = 60 * time.Second

type Peer struct {
	ID        string    `json:"id"`
	StartedAt time.Time `json:"started_at"`

	closed    bool
	log       *logger.Logger
	pingTimer *time.Timer
	reader    *bufio.Reader
	conn      net.Conn
}

func New(conn net.Conn) *Peer {
	id := uuid.NewString()
	reader := bufio.NewReader(conn)
	now := time.Now()
	log := logger.New(fmt.Sprintf("vayeate:peer:%s", id))
	self := Peer{id, now, false, log, nil, reader, conn}
	self.pingTimer = time.AfterFunc(timeout, onTimeout(&self))
	return &self
}

func (self *Peer) Close() {
	self.closed = true
	self.conn.Close()
	self.pingTimer.Stop()
}

func (self *Peer) Closed() bool {
	return self.closed
}

func (self *Peer) GetRemoteAddress() string {
	return self.conn.RemoteAddr().String()
}

func (self *Peer) Read() (*Message, error) {
	f, err := DecodeMessage(self.reader)

	if err == io.EOF {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	self.pingTimer.Reset(timeout)
	return f, nil
}

func (self *Peer) Write(m *Message) error {
	_, err := self.conn.Write(m.Encode())
	return err
}

func onTimeout(self *Peer) func() {
	return func() {
		self.Close()
		self.log.Infoln("closed due to inactivity")
	}
}
