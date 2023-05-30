package peer

import (
	"bufio"
	"bytes"
	"encoding/gob"
)

type OpCode uint8

const (
	OPEN         OpCode = 0
	OPEN_SUCCESS OpCode = 1
	PING         OpCode = 2
	PONG         OpCode = 3
	ASSERT       OpCode = 4
	PRODUCE      OpCode = 5
	ACK          OpCode = 6
)

type MessageHeaders struct {
	FromID   string
	Username string
	Password string
}

type Message struct {
	Code    OpCode
	Headers MessageHeaders
	Subject []byte
	Body    []byte
}

func NewOpenMessage(id string, username string, password string) *Message {
	return &Message{
		OPEN,
		MessageHeaders{
			id,
			username,
			password,
		},
		nil,
		nil,
	}
}

func NewOpenSuccessMessage(id string, username string, password string) *Message {
	return &Message{
		OPEN_SUCCESS,
		MessageHeaders{
			id,
			username,
			password,
		},
		nil,
		nil,
	}
}

func NewPingMessage(id string, username string, password string) *Message {
	return &Message{
		PING,
		MessageHeaders{
			id,
			username,
			password,
		},
		nil,
		nil,
	}
}

func NewPongMessage(id string, username string, password string) *Message {
	return &Message{
		PONG,
		MessageHeaders{
			id,
			username,
			password,
		},
		nil,
		nil,
	}
}

func DecodeMessage(reader *bufio.Reader) (*Message, error) {
	var self Message
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&self)

	if err != nil {
		return nil, err
	}

	return &self, nil
}

func (self *Message) Encode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(self)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (self *Message) GetSubject() string {
	return string(self.Subject)
}

func (self *Message) GetBody() string {
	return string(self.Body)
}

func (self *Message) IsOpen() bool {
	return self.Code == OPEN
}

func (self *Message) IsOpenSuccess() bool {
	return self.Code == OPEN_SUCCESS
}

func (self *Message) IsPing() bool {
	return self.Code == PING
}

func (self *Message) IsPong() bool {
	return self.Code == PONG
}

func (self *Message) IsAssert() bool {
	return self.Code == ASSERT
}

func (self *Message) IsProduce() bool {
	return self.Code == PRODUCE
}

func (self *Message) IsAck() bool {
	return self.Code == ACK
}
