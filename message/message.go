package message

import (
	"encoding/binary"
	"io"
	"log"
)

type messageID uint8

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte
}

func (m *Message) Serialize() []byte {
	if m == nil {
		log.Println("got empty message!")
		return make([]byte, 4)
	}

	length := len(m.Payload) + 1
	buf := make([]byte, length+4)
	binary.BigEndian.PutUint32(buf[0:4], uint32(length))
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

func Read(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		log.Println("error reading message length into length buf!")
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	if length == 0 {
		log.Println("received message has a length of 0")
		return nil, err
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		log.Println("error reading message into message buf!")
		return nil, err
	}

	return &Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}, nil
}
