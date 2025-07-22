package client

import (
	"log"
	"net"
	"time"

	"github.com/aryanbroy/bittorrent/peers"
)

type Client struct {
	Conn   net.Conn
	Choked bool
}

func New(peer peers.Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		log.Println("Unable to establish tcp connection...")
		return nil, err
	}
}
