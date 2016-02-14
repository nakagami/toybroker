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

type MemoryHook struct {
	clientMap        map[string]Client
	messageBufferMap map[string]MemoryMessageBuffer
	clientMapMutex   sync.Mutex
}

func NewMemoryHook() MemoryHook {
	return MemoryHook{
		clientMap:        make(map[string]Client),
		messageBufferMap: make(map[string]MemoryMessageBuffer),
	}
}

func (h MemoryHook) Login(conn net.Conn, clientID string, loginName string, loginPassword string) byte {
	return CONNACK_Success
}

func (h MemoryHook) Logout(clientID string) {
}

func (h MemoryHook) GetClient(clientID string) Client {
	h.clientMapMutex.Lock()
	defer h.clientMapMutex.Unlock()
	return h.clientMap[clientID]
}

func (h MemoryHook) SetClient(clientID string, loginName string, conn net.Conn) Client {
	client := NewClient(clientID, loginName, conn)
	h.clientMapMutex.Lock()
	defer h.clientMapMutex.Unlock()
	h.clientMap[client.GetClientID()] = client
	return client
}

func (h MemoryHook) GetMessageBuffer(clientID string) MessageBuffer {
	h.clientMapMutex.Lock()
	defer h.clientMapMutex.Unlock()
	mb, ok := h.messageBufferMap[clientID]
	if !ok {
		mb := NewMemoryMessageBuffer(clientID, 10)
		h.messageBufferMap[clientID] = mb
	}
	return mb
}

func (h MemoryHook) Subscribe(topics Topics, topicName string, clientID string, qos int) int {
	topics.Add(topicName, clientID, qos)
	return 0 // QoS
}

func (h MemoryHook) Unsubscribe(topics Topics, topicName string, clientID string) {
	topics.Remove(topicName, clientID)
	return
}
