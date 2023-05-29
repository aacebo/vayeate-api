package node

import (
	"fmt"
	"net"
	"strconv"
	"vayeate-api/logger"
	"vayeate-api/sync"

	"github.com/google/uuid"
)

type Node struct {
	ID   string
	Port int

	log      *logger.Logger
	listener net.Listener
	sockets  sync.SyncMap[string, *Socket]
}

func New(port string) (*Node, error) {
	id := uuid.NewString()
	p, err := strconv.Atoi(port)

	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", p))

	if err != nil {
		return nil, err
	}

	self := Node{
		ID:       id,
		Port:     p,
		log:      logger.New(fmt.Sprintf("vayeate:node:%s", id)),
		listener: listener,
		sockets:  sync.NewSyncMap[string, *Socket](),
	}

	return &self, nil
}

func (self *Node) Listen() error {
	for {
		conn, err := self.listener.Accept()

		if err != nil {
			return err
		}

		go self.onConnection(conn)
	}
}

func (self *Node) Close() {
	self.listener.Close()
}

func (self *Node) GetSockets() []*Socket {
	return self.sockets.Slice()
}

func (self *Node) onConnection(conn net.Conn) {
	s := NewSocket(conn)
	self.sockets.Set(s.ID, s)

	defer func() {
		s.Close()
		self.sockets.Del(s.ID)
	}()

	for {
		if s.closed == true {
			return
		}

		m, err := s.Read()

		if m == nil || err != nil {
			if err != nil {
				self.log.Warn(err)

				if err == InvalidFormatError {
					return
				}
			}

			continue
		}

		if m.IsClose() {
			return
		} else if m.IsPing() {
			self.onPing(s)
		} else {
			self.onMessage(s, m)
		}
	}
}

func (self *Node) onPing(s *Socket) {
	err := s.Write(NewPongMessage())

	if err != nil {
		self.log.Warn(err)
	}
}

func (self *Node) onMessage(s *Socket, m *Message) {

}
