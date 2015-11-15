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
	"sync"
)

type MemoryMessageBuffer struct {
	m map[uint16][]byte
	sync.RWMutex
}

func NewMemoryMessageBuffer(clientID string) MemoryMessageBuffer {
	return MemoryMessageBuffer{
		m: make(map[uint16][]byte),
	}
}

func (m MemoryMessageBuffer) Set(messageID uint16, payload []byte) {
	m.Lock()
	defer m.Unlock()
	m.m[messageID] = payload
}

func (m MemoryMessageBuffer) Get(messageID uint16) []byte {
	m.RLock()
	defer m.RUnlock()
	payload, ok := m.m[messageID]
	if !ok {
		return nil
	}
	return payload
}

func (m MemoryMessageBuffer) Delete(messageID uint16) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, messageID)
}

func (m MemoryMessageBuffer) List() []uint16 {
	m.RLock()
	defer m.RUnlock()
	list := make([]uint16, 0)
	for k, _ := range m.m {
		list = append(list, k)
	}
	return list
}
