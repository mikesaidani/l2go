package models

import (
	"net"
)

type Client struct {
	SessionID []byte
	Socket    net.Conn
}

func NewClient() *Client {
  return &Client{}
}
