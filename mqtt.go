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
	"fmt"
)

func MqttMainLoop(conn net.Conn) {
	command, _, remaining, err := readMessage(conn)
	if command != CONNECT || err != nil {
		sendToClientConn(packCONNACK(CONNACK_Rejected))
		return
	}
	clientID, _, _, loginName, loginPassword, err := unpackCONNECT(remaining)
	status := login(clientID, loginName, loginPassword)
    sendToClientID(packCONNACK(status), clientID)

	fmt.Println("clientID=", clientID)

	for {
		command, _, remaining, err = readMessage()
		if err != nil {
			logout(clientID)
			break
		}
		switch command {
		case PUBLISH:
			fmt.Println("PUBLISH")
			topic, messageID, payload, err := unpackPUBLISH(remaining)
			fmt.Println(topic, messageID, payload, err)
			target_clients := messages.CreateMessages(clientID, messageID, topics.GetClients(topic), payload)

			for _, c := range target_clients {
				peer := clients.GetPeer(c.ClientID)
				if peer != nil {
					peer.WriteChan <- packPUBLISH(topic, 1, payload)
				}
			}
		case PUBACK:
			fmt.Println("PUBACK")
			messageID, _ := unpackPUBACK(remaining)
			ackClientID, ackMessageID := messages.RecvPUBACK(clientID, messageID)
			if ackMessageID != 0 {
				ackPeer := clients.GetPeer(ackClientID)
				if ackPeer != nil {
					ackPeer.WriteChan <- packPUBACK(ackMessageID)
				}
			}
		case PUBREL:
			fmt.Println("PUBREL")
			messageID, _ := unpackPUBREL(remaining)
			peer.WriteChan <- packPUBCOMP(messageID)
		case SUBSCRIBE:
			fmt.Println("SUBSCRIBE")
			messageID, subscribe_topics, err := unpackSUBSCRIBE(remaining)
			fmt.Println(messageID, topics, err)
			qos := make([]byte, len(subscribe_topics))
			for i, topic := range subscribe_topics {
				qos[i] = topics.Subscribe(topic.TopicID, clientID, topic.Qos)
			}
			peer.WriteChan <- packSUBACK(messageID, qos)
		case UNSUBSCRIBE:
			fmt.Println("UNSUBSCRIBE")
			messageID, unsubscribe_topics, err := unpackUNSUBSCRIBE(remaining)
			fmt.Println(messageID, unsubscribe_topics, err)
			for _, topic := range unsubscribe_topics {
				topics.Unsubscribe(topic, clientID)
			}
		case PINGREQ:
			fmt.Println("PINGREQ")
			peer.WriteChan <- packPINGRESP()
		case DISCONNECT:
			fmt.Println("DISCONNECT")
			clients.Close(clientID, peer)
			break
		default:
			fmt.Println("Invalid Command")
			clients.Close(clientID, peer)
			break
		}
	}
}
