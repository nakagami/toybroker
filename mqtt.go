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
	"fmt"
	"net"
)

func MqttMainLoop(conn net.Conn, topics Topics, hook Hook) {
	connect, _, _, _, _, remaining, err := readMessage(conn)
	if connect != CONNECT || err != nil {
		conn.Write([]byte{CONNACK * 16, 2, 0, CONNACK_Rejected})
		return
	}
	clientID, _, _, loginName, loginPassword, err := unpackCONNECT(remaining)
	status, client := hook.Login(conn, clientID, loginName, loginPassword)
	if client == nill {
		conn.Write([]byte{CONNACK * 16, 2, 0, status})
		return
	}

	client.Send(packCONNACK(status))

	for {
		command, _, header_qos, retain, _, remaining, err := readMessage(conn)
		if err != nil {
			hook.Logout(clientID)
			break
		}
		switch command {
		case PUBLISH:
			topic, messageID, payload, err := unpackPUBLISH(remaining)
			debugOutput(fmt.Sprintf("PUBLISH:%s,%d,%v,%v", topic, messageID, payload, err))

			clientList, qosList := topics.List(topic)
			for i, clientID := range clientList {
				target := hook.GetClient(clientID)
				if qosList[i] > 0 && header_qos > 0 {
					qosList[i] = 1
				}
				target.Publish(false, qosList[i], false, topic, payload)
			}
			if retain {
				topics.AddRetainMessage(topic, payload)
			}
		case PUBACK:
			debugOutput("PUBACK")
		case PUBREL:
			debugOutput("PUBREL")
		case SUBSCRIBE:
			messageID, subscribe_topics, err := unpackSUBSCRIBE(remaining)
			debugOutput(fmt.Sprintf("SUBSCRIBE:%v,%d,%v", subscribe_topics, messageID, err))
			qos := make([]byte, len(subscribe_topics))
			for i, topicName := range subscribe_topics {
				qos[i] = byte(hook.Subscribe(topics, topicName, clientID, 0))
			}
			client.Send(packSUBACK(messageID, qos))

			// Publish retain messages
			for i, topicName := range subscribe_topics {
				payload := topics.GetRetainMessage(topicName)
				debugOutput(fmt.Sprintf("GetRetainMessage(%s)=%v", topicName, payload))

				if len(payload) != 0 {
					debugOutput(fmt.Sprintf("SendRetainMessage:%s->%v", topicName, payload))
					client.Publish(false, int(qos[i]), true, topicName, payload)
				}
			}

		case UNSUBSCRIBE:
			messageID, unsubscribe_topics, err := unpackUNSUBSCRIBE(remaining)
			debugOutput(fmt.Sprintf("UNSUBSCRIBE:%v,%d,%v", unsubscribe_topics, messageID, err))
			for _, topicName := range unsubscribe_topics {
				hook.Unsubscribe(topics, topicName, clientID)
			}
			client.Send(packUNSUBACK(messageID))
		case PINGREQ:
			debugOutput("PINGREQ")
			client.Send(packPINGRESP())
		case DISCONNECT:
			debugOutput("DISCONNECT")
			hook.Logout(clientID)
			break
		default:
			debugOutput("Invalid Command")
			hook.Logout(clientID)
			break
		}
	}
}
