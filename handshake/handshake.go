package handshake

import (
	"io"
	"log"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buf[curr:], []byte(h.Pstr))
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

func Read(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		log.Println("error reading into length buffer")
		return nil, err
	}

	pstrLen := int(lengthBuf[0])

	if pstrLen == 0 {
		log.Println("protocol identifier length cannot be 0")
		return nil, err
	}

	handshakeBuf := make([]byte, pstrLen+48)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		log.Println("error reading into handshake buffer")
		return nil, err
	}

	pstr := handshakeBuf[:pstrLen]

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakeBuf[pstrLen+8:pstrLen+8+20])
	copy(peerID[:], handshakeBuf[pstrLen+8+20:])

	return &Handshake{
		Pstr:     string(pstr),
		InfoHash: infoHash,
		PeerID:   peerID,
	}, nil
}
