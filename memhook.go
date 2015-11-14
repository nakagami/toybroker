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
	"net"
	"sync"
)

type MemoryHook struct {
    clientMap map[string]*Client
    clientMapMutex sync.Mutex
}

func NewMemoryHook() * MemoryHook {
    return &MemoryHook{
        clientMap: make(map[string]*Client),
    }
}

func (h *MemoryHook) Login(conn net.Conn, clientID string, loginName string, loginPassword string) byte {
	return CONNACK_Success
}

func (h *MemoryHook) Logout(clientID string) {
}

func (h *MemoryHook) GetClient(clientID string) *Client {
	clientMapMutex.Lock()
	defer clientMapMutex.Unlock()
	return clientMap[clientID]
}

func (h *MemoryHook) SetClient(client *Client) {
	clientMapMutex.Lock()
	defer clientMapMutex.Unlock()
	clientMap[client.GetClientID()] = client
}

func (h *MemoryHook) Subscribe(topics *Topics, topicName string, clientID string) byte {
	topics.Add(topicName, clientID)
	return 0 // QoS
}

func (h *MemoryHook) Unsubscribe(topics *Topics, topicName string, clientID string) {
	topics.Remove(topicName, clientID)
	return
}
