package models

import (
	"crypto/rand"
	"net"
)

type Client struct {
	Account   Account
	SessionID []byte
	Socket    net.Conn
}

func NewClient() *Client {
	id := make([]byte, 16)
	_, err := rand.Read(id)

	if err != nil {
		return nil
	}
  return &Client{SessionID: id}
}
