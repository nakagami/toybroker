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
	"fmt"
	"net"
)

func MqttMainLoop(conn net.Conn) {
	command, _, remaining, err := readMessage(conn)
	if command != CONNECT || err != nil {
		sendToConn(packCONNACK(CONNACK_Rejected), conn)
		return
	}
	clientID, _, _, loginName, loginPassword, err := unpackCONNECT(remaining)
	status := login(conn, clientID, loginName, loginPassword)
	client := NewGetClient(clientID, conn, 0)
	setClient(client)

	sendToConn(packCONNACK(status), conn)
	if status != CONNACK_Success {
		return
	}

	for {
		command, _, remaining, err = readMessage(conn)
		if err != nil {
			logout(clientID)
			break
		}
		switch command {
		case PUBLISH:
			topic, messageID, payload, err := unpackPUBLISH(remaining)
			debugOutput(fmt.Sprintf("PUBLISH:%s,%d,%v,%v", topic, messageID, payload, err))

			for _, clientID := range getClientListByTopic(topic) {
				messageID = getNextMessageID(clientID)
				data := packPUBLISH(topic, messageID, payload)
				sendToClient(data, clientID)
			}
		case PUBACK:
			debugOutput("PUBACK")
		case PUBREL:
			debugOutput("PUBREL")
		case SUBSCRIBE:
			messageID, subscribe_topics, err := unpackSUBSCRIBE(remaining)
			debugOutput(fmt.Sprintf("SUBSCRIBE:%v,%d,%v", subscribe_topics, messageID, err))
			qos := make([]byte, len(subscribe_topics))
			for i, topic := range subscribe_topics {
				qos[i] = subscribe(topic, clientID)
			}
			sendToConn(packSUBACK(messageID, qos), conn)
		case UNSUBSCRIBE:
			messageID, unsubscribe_topics, err := unpackUNSUBSCRIBE(remaining)
			debugOutput(fmt.Sprintf("UNSUBSCRIBE:%v,%d,%v", unsubscribe_topics, messageID, err))
			for _, topic := range unsubscribe_topics {
				unsubscribe(topic, clientID)
			}
			sendToConn(packUNSUBACK(messageID), conn)
		case PINGREQ:
			debugOutput("PINGREQ")
			sendToConn(packPINGRESP(), conn)
		case DISCONNECT:
			debugOutput("DISCONNECT")
			logout(clientID)
			break
		default:
			debugOutput("Invalid Command")
			logout(clientID)
			break
		}
	}
}
