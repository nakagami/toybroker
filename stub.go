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

package main

import (
	"github.com/streamrail/concurrent-map"
	"net"
	"sync"
)

var clientMap = cmap.New()
var topicMap = cmap.New()

// ------------------------------ Client --------------------------------------

type Client struct {
	clientID      string
	conn          net.Conn
	nextMessageID uint16
	sync.RWMutex
}

// ------------------------------ Topic ---------------------------------------

type Topic struct {
	Name string
	m    map[string]bool
	sync.RWMutex
}

func NewTopic(name string) *Topic {
	return &Topic{
		Name: name,
		m:    make(map[string]bool),
	}
}

func (s *Topic) Add(name string) {
	s.Lock()
	defer s.Unlock()
	s.m[name] = true
}

func (s *Topic) Remove(name string) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, name)
}

func (s *Topic) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[string]bool)
}

func (s *Topic) List() []string {
	s.RLock()
	defer s.RUnlock()
	list := make([]string, 0)
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

// -----------------------------------------------------------------------

func initialize_stub() {
}

func login(clientID string, loginName string, loginPassword string) byte {
	return CONNACK_Success
}

func logout(clientID string) {
}

func getNextMessageID(clientID string) uint16 {
	return 1
}

func getClientListByTopic(topicName string) []string {
	return []string{}
}

func sendToClient(data []byte, clientID string) bool {
	return true
}

func sendToConn(data []byte, conn net.Conn) {
	conn.Write(data)
}

func subscribe(topicID string, clientID string) byte {
	return 0 // QoS
}

func unsubscribe(topicID string, clientID string) {
	return
}
