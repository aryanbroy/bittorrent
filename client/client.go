package client

import (
	"log"
	"net"
	"time"

	"github.com/aryanbroy/bittorrent/bitfield"
	"github.com/aryanbroy/bittorrent/message"
	"github.com/aryanbroy/bittorrent/peers"
)

type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield bitfield.Bitfield
}

func (c *Client) SendHave(index int) error {
	msg := message.FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func New(peer peers.Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		log.Println("Unable to establish tcp connection...")
		return nil, err
	}
}

func (c *Client) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) Read() (*message.Message, error) {
	msg, err := message.Read(c.Conn)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Client) SendRequest(index, begin, length int) error {
	req := message.FormatRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}
