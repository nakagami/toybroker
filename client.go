/*******************************************************************************
The MIT License (MIT)

Copyright (c) 2015 Hajime Nakagami

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*******************************************************************************/

package toybroker

import (
	"net"
	"sync"
)

type Client struct {
	clientID         string
	loginName        string
	conn             net.Conn
	currentMessageID uint16
	sync.RWMutex
}

func NewClient(id string, name string, c net.Conn, oldClient *Client) *Client {
	client := &Client{
		clientID:  id,
		loginName: name,
		conn:      c,
	}
	return client
}

func (c *Client) GetClientID() string {
	return c.clientID
}

func (c *Client) GetLoginName() string {
	return c.loginName
}

func (c *Client) GetConn() net.Conn {
	return c.conn
}

func (c *Client) getNextMessageID() uint16 {
	c.Lock()
	defer c.Unlock()
	c.currentMessageID++
	if c.currentMessageID == 0 {
		c.currentMessageID++
	}
	return c.currentMessageID
}

func (c *Client) Publish(dup bool, qos int, retain bool, topic string, payload []byte) {
	c.conn.Write(packPUBLISH(dup, qos, retain, topic, c.getNextMessageID(), payload))
}

func (c *Client) Send(data []byte) {
	c.conn.Write(data)
}
