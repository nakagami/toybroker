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
	"sort"
	"sync"
)

type MemoryTopics struct {
	mapClientQoS map[string]map[string]int
	mapRetain    map[string][]byte
	sync.RWMutex
}

func NewMemoryTopics() MemoryTopics {
	return MemoryTopics{
		mapClientQoS: make(map[string]map[string]int),
		mapRetain:    make(map[string][]byte),
	}
}

func (t MemoryTopics) Add(topicName string, clientID string, qos int) {
	t.Lock()
	defer t.Unlock()
	topic, ok := t.mapClientQoS[topicName]
	if !ok {
		topic = make(map[string]int)
		t.mapClientQoS[topicName] = topic
	}
	topic[clientID] = qos
}

func (t MemoryTopics) Remove(topicName string, clientID string) {
	t.Lock()
	defer t.Unlock()
	topic, ok := t.mapClientQoS[topicName]
	if ok {
		delete(topic, clientID)
	}
}

func (t MemoryTopics) TopicList() []string {
	t.RLock()
	defer t.RUnlock()
	list := make([]string, 0)
	for k, _ := range t.mapClientQoS {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

func (t MemoryTopics) List(topicName string) ([]string, []int) {
	t.RLock()
	defer t.RUnlock()
	clientList := make([]string, 0)
	qosList := make([]int, 0)
	topic, ok := t.mapClientQoS[topicName]
	if ok {
		for c := range topic {
			clientList = append(clientList, c)
			qosList = append(qosList, topic[c])
		}
	}
	return clientList, qosList
}

func (t MemoryTopics) AddRetainMessage(topicName string, payload []byte) {
	t.Lock()
	defer t.Unlock()
	t.mapRetain[topicName] = payload
}

func (t MemoryTopics) GetRetainMessage(topicName string) []byte {
	t.RLock()
	defer t.RUnlock()
	retain, ok := t.mapRetain[topicName]
	if ok {
		return retain
	}
	return []byte{}
}
