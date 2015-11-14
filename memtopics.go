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
	"sort"
	"sync"
)

type Topics struct {
	m map[string]map[string]bool
	sync.RWMutex
}

func NewTopics() Topics {
	return Topics{
		m: make(map[string]map[string]bool),
	}
}

func (t Topics) Add(topicName string, clientID string) {
	t.Lock()
	defer t.Unlock()
	topic, ok := t.m[topicName]
	if !ok {
		topic = make(map[string]bool)
		t.m[topicName] = topic
	}
	topic[clientID] = true
}

func (t Topics) Remove(topicName string, clientID string) {
	t.Lock()
	defer t.Unlock()
	topic, ok := t.m[topicName]
	if ok {
		delete(topic, clientID)
	}
}

func (t Topics) TopicList() []string {
	t.RLock()
	defer t.RUnlock()
	list := make([]string, 0)
	for k, _ := range t.m {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

func (t Topics) List(topicName string) []string {
	t.RLock()
	defer t.RUnlock()
	list := make([]string, 0)
	topic, ok := t.m[topicName]
	if ok {
		for item := range topic {
			list = append(list, item)
		}
	}
	return list
}
