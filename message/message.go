package message

import (
	"encoding/binary"
	"fmt"
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

func ParseHave(msg *Message) (int, error) {
	if msg.ID != MsgHave {
		return 0, fmt.Errorf("expected id to be %v, got %v instead", MsgHave, msg.ID)
	}

	if len(msg.Payload) < 4 {
		return 0, fmt.Errorf("expected payload length to be %v bytes, got %v bytes instead", 4, len(msg.Payload))
	}

	index := binary.BigEndian.Uint32(msg.Payload)
	return int(index), nil
}

func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID != MsgPiece {
		return 0, fmt.Errorf("expected id to be %v, got %v instead", MsgPiece, msg.ID)
	}

	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("payload too short. %d < 8", len(msg.Payload))
	}

	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("expected index %d, %d != %d", index, parsedIndex, index)
	}

	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(msg.Payload) {
		return 0, fmt.Errorf("data begin offset too large. %d >= %d", begin, len(msg.Payload))
	}

	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf("data too long [%d], with offset %d and length %d", len(data), begin, len(buf))
	}

	copy(buf[begin:], data)
	return len(data), nil
}

func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: MsgRequest, Payload: payload}
}

func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{ID: MsgHave, Payload: payload}
}
