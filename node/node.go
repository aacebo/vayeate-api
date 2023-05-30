package node

import (
	"fmt"
	"net"
	"strconv"
	"time"
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

	username       string
	password       string
	entryAddress   string
	log            *logger.Logger
	socketListener net.Listener
	peerListener   net.Listener
	sockets        sync.SyncMap[string, *socket.Socket]
	peers          sync.SyncMap[string, *peer.Peer]
}

func New(socketPort string, peerPort string, username string, password string, entryAddress string) (*Node, error) {
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
		username:       username,
		password:       password,
		entryAddress:   entryAddress,
		log:            logger.New(fmt.Sprintf("vayeate:node:%s", id)),
		socketListener: socketListener,
		peerListener:   peerListener,
		sockets:        sync.NewSyncMap[string, *socket.Socket](),
		peers:          sync.NewSyncMap[string, *peer.Peer](),
	}

	if entryAddress != "" {
		// connect to other nodes
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
		} else {
			err = self.onSocketMessage(s, m)

			if err != nil {
				return
			}
		}
	}
}

func (self *Node) onPeerConnection(conn net.Conn) {
	p, err := peer.FromConnection(self.ID, self.username, self.password, conn)

	if err != nil {
		self.log.Warn(err)
		return
	}

	self.peers.Set(p.ID, p)
	pingTimer := time.AfterFunc(30*time.Second, func() {
		err := p.Write(peer.NewPingMessage(p.ID, self.username, self.password))

		if err != nil {
			self.log.Warn(err)
		}
	})

	defer func() {
		p.Close()
		pingTimer.Stop()
		self.peers.Del(p.ID)
	}()

	for {
		if p.Closed() {
			return
		}

		m, err := p.Read()

		if m == nil || err != nil {
			if err != nil {
				self.log.Warn(err)

				if err == peer.UnauthorizedError {
					return
				}
			}

			continue
		}

		err = self.onPeerMessage(p, m)

		if err != nil {
			return
		}
	}
}

func (self *Node) onSocketMessage(s *socket.Socket, m *socket.Message) error {
	var err error

	if m.IsPing() {
		err = s.Write(socket.NewPongMessage())
	}

	return err
}

func (self *Node) onPeerMessage(p *peer.Peer, m *peer.Message) error {
	var err error

	if m.IsPing() {
		err = p.Write(peer.NewPongMessage(self.ID, self.username, self.password))
	}

	return err
}
