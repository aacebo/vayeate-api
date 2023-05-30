package socket

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type OpCode uint8

const (
	CLOSE    OpCode = 0 // <code::>
	PING     OpCode = 1 // <code::>
	PONG     OpCode = 2 // <code::>
	ASSERT   OpCode = 3 // <code:subject:>
	PRODUCE  OpCode = 4 // <code:subject:body>
	CONSUME  OpCode = 5 // <code:subject:>
	ACK      OpCode = 6 // <code:subject:>
	DELEGATE OpCode = 7 // <code:subject:body>
)

const (
	START     = byte('<') // start byte
	END       = byte('>') // end byte
	DELIMITER = byte(':') // slice delimiter
)

var InvalidFormatError = errors.New("socket: invalid message format")

// A Message comprised of many bytes, utilizing a custom Netstring format
// ex. <code:subject:body>
// https://cr.yp.to/proto/netstrings.txt
type Message struct {
	Code    OpCode
	Subject []byte
	Body    []byte
}

func NewMessage(code OpCode, subject []byte, body []byte) *Message {
	return &Message{code, subject, body}
}

func NewCloseMessage() *Message {
	return &Message{CLOSE, []byte{}, []byte{}}
}

func NewPingMessage() *Message {
	return &Message{PING, []byte{}, []byte{}}
}

func NewPongMessage() *Message {
	return &Message{PONG, []byte{}, []byte{}}
}

func NewAssertMessage(subject []byte) *Message {
	return &Message{ASSERT, subject, []byte{}}
}

func NewProduceMessage(subject []byte, body []byte) *Message {
	return &Message{PRODUCE, subject, body}
}

func NewConsumeMessage(subject []byte) *Message {
	return &Message{CONSUME, subject, []byte{}}
}

func NewAckMessage(subject []byte) *Message {
	return &Message{ACK, subject, []byte{}}
}

func NewDelegateMessage(subject []byte, body []byte) *Message {
	return &Message{DELEGATE, subject, body}
}

func DecodeMessage(reader *bufio.Reader) (*Message, error) {
	subject := []byte{}
	body := []byte{}
	b, err := reader.ReadByte()

	if err != nil {
		return nil, err
	}

	if b != START {
		return nil, InvalidFormatError
	}

	// read opcode
	b, err = reader.ReadByte()

	if err != nil {
		return nil, err
	}

	t, err := strconv.Atoi(string(b))

	if err != nil {
		return nil, err
	}

	code := OpCode(t)
	b, err = reader.ReadByte()

	if err != nil {
		return nil, err
	}

	if b != DELIMITER {
		return nil, InvalidFormatError
	}

	// read subject
	for {
		b, err := reader.ReadByte()

		if err == io.EOF {
			return nil, InvalidFormatError
		}

		if err != nil {
			return nil, err
		}

		if b == DELIMITER {
			break
		}

		subject = append(subject, b)
	}

	// read body
	for {
		b, err := reader.ReadByte()

		if err == io.EOF {
			return nil, InvalidFormatError
		}

		if err != nil {
			return nil, err
		}

		if b == END {
			break
		}

		body = append(body, b)
	}

	return &Message{code, subject, body}, nil
}

func (self *Message) Encode() []byte {
	data := []byte{}
	code := []byte(strconv.Itoa(int(self.Code)))

	data = append(data, START)
	data = append(data, code...)
	data = append(data, DELIMITER)
	data = append(data, self.Subject...)
	data = append(data, DELIMITER)
	data = append(data, self.Body...)
	data = append(data, END)

	return data
}

func (self *Message) GetSubject() string {
	return string(self.Subject)
}

func (self *Message) GetBody() string {
	return string(self.Body)
}

func (self *Message) IsClose() bool {
	return self.Code == CLOSE
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

func (self *Message) IsConsume() bool {
	return self.Code == CONSUME
}

func (self *Message) IsAck() bool {
	return self.Code == ACK
}

func (self *Message) IsDelegate() bool {
	return self.Code == DELEGATE
}
