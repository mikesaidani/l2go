package models

import (
	"net"
)

type Client struct {
	SessionID []byte
	Socket    net.Conn
  Cipher    XorCipher
}

func NewClient() *Client {
  return &Client{Cipher: NewXorCipher()}
}
