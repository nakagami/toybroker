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

import "strings"

type Topics interface {
	// Add client and it's QoS to topic
	Add(topicName string, clientID string, qos int)
	// Remove client (and it's QoS) from topic
	Remove(topicName string, clientID string)
	// Get all topic names
	TopicList() []string
	// Get Client IDs and their QoS by topic name
	List(topicName string) ([]string, []int)

	// Add (or replace) retain message to topic
	AddRetainMessage(topicName string, payload []byte)
	// Get retain message from topic
	GetRetainMessage(topicName string) []byte
}

func HasTopicWildCard(pat string) bool {
	return strings.Contains("#", pat) || strings.Contains("+", pat)
}

func TopicMatchList(pat string, topicList []string) []string {
	// TODO
	matchList := make([]string, 0, len(topicList))

	for _, topicName := range topicList {
		if pat == topicName {
			matchList = append(matchList, topicName)
		}
	}

	return matchList
}
