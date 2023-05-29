package node

import (
	"fmt"
	"net"
	"strconv"
	"vayeate-api/logger"
	"vayeate-api/peer"
	"vayeate-api/socket"
	"vayeate-api/sync"

	"github.com/google/uuid"
)

type Node struct {
	ID         string
	SocketPort int
	PeerPort   int

	log            *logger.Logger
	socketListener net.Listener
	peerListener   net.Listener
	sockets        sync.SyncMap[string, *socket.Socket]
	peers          sync.SyncMap[string, *peer.Peer]
}

func New(socketPort string, peerPort string) (*Node, error) {
	id := uuid.NewString()
	sp, err := strconv.Atoi(socketPort)

	if err != nil {
		return nil, err
	}

	socketListener, err := net.Listen("tcp", fmt.Sprintf(":%d", sp))

	if err != nil {
		return nil, err
	}

	pp, err := strconv.Atoi(peerPort)

	if err != nil {
		return nil, err
	}

	peerListener, err := net.Listen("tcp", fmt.Sprintf(":%d", pp))

	self := Node{
		ID:             id,
		SocketPort:     sp,
		PeerPort:       pp,
		log:            logger.New(fmt.Sprintf("vayeate:node:%s", id)),
		socketListener: socketListener,
		peerListener:   peerListener,
		sockets:        sync.NewSyncMap[string, *socket.Socket](),
		peers:          sync.NewSyncMap[string, *peer.Peer](),
	}

	return &self, nil
}

func (self *Node) Listen() {
	go self.listenSockets()
	go self.listenPeers()
}

func (self *Node) Close() {
	self.socketListener.Close()
	self.peerListener.Close()
}

func (self *Node) GetSockets() []*socket.Socket {
	return self.sockets.Slice()
}

func (self *Node) GetPeers() []*peer.Peer {
	return self.peers.Slice()
}

func (self *Node) listenSockets() error {
	for {
		conn, err := self.socketListener.Accept()

		if err != nil {
			return err
		}

		go self.onSocketConnection(conn)
	}
}

func (self *Node) listenPeers() error {
	for {
		conn, err := self.peerListener.Accept()

		if err != nil {
			return err
		}

		go self.onPeerConnection(conn)
	}
}

func (self *Node) onSocketConnection(conn net.Conn) {
	s := socket.New(conn)
	self.sockets.Set(s.ID, s)

	defer func() {
		s.Close()
		self.sockets.Del(s.ID)
	}()

	for {
		if s.Closed() {
			return
		}

		m, err := s.Read()

		if m == nil || err != nil {
			if err != nil {
				self.log.Warn(err)

				if err == socket.InvalidFormatError {
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
			self.onSocketMessage(s, m)
		}
	}
}

func (self *Node) onPeerConnection(conn net.Conn) {
}

func (self *Node) onPing(s *socket.Socket) {
	err := s.Write(socket.NewPongMessage())

	if err != nil {
		self.log.Warn(err)
	}
}

func (self *Node) onSocketMessage(s *socket.Socket, m *socket.Message) {

}
