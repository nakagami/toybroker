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
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

func bytes_to_uint16(b []byte) uint16 {
	var val uint16
	buffer := bytes.NewBuffer(b)
	binary.Read(buffer, binary.BigEndian, &val)
	return val
}

func bytes_to_str(b []byte) string {
	return bytes.NewBuffer(b).String()
}

func str_to_bytes(s string) []byte {
	b := bytes.NewBufferString(s).Bytes()
	ln := len(b)
	return bytes.Join([][]byte{[]byte{byte(ln / 256), byte(ln)}, b}, nil)
}

func encodeRemainLength(x int) []byte {
	var bs []byte
	if x == 0 {
		bs = make([]byte, 1)
		bs[0] = 0
	} else {
		bs = make([]byte, 4)
		var i int
		for i = 0; x > 0; i++ {
			digit := byte(x % 128)
			x = x / 128
			if x > 0 {
				digit = digit | 0x80
			}
			bs[i] = digit
		}
		bs = bs[0:i]
	}

	return bs
}

func decodeRemainLength(bs []byte) int {
	var multiplier = 1
	var value = 0

	for _, b := range bs {
		value += int(b&127) * multiplier
		multiplier *= 128
	}

	return value
}

/* Pack return data */
func packPUBLISH(topic string, messageID uint16, payload []byte) []byte {
	remaining := bytes.Join([][]byte{
		str_to_bytes(topic),
		[]byte{byte(messageID / 256), byte(messageID % 256)},
		payload,
	}, nil)
	return bytes.Join([][]byte{
		[]byte{PUBLISH * 16},
		encodeRemainLength(len(remaining)),
		remaining,
	}, nil)
}

func packCONNACK(returnCode byte) []byte {
	return []byte{CONNACK * 16, 2, 0, returnCode}
}

func packPUBACK(messageID uint16) []byte {
	return []byte{CONNACK * 16, 2, byte(messageID / 256), byte(messageID % 256)}
}

func packPUBREC(messageID uint16) []byte {
	return []byte{PUBREC * 16, 2, byte(messageID / 256), byte(messageID % 256)}
}

func packPUBCOMP(messageID uint16) []byte {
	return []byte{PUBCOMP * 16, 2, byte(messageID / 256), byte(messageID % 256)}
}

func packSUBACK(messageID uint16, qos []byte) []byte {
	ln := len(qos) + 2
	return bytes.Join([][]byte{
		[]byte{SUBACK * 16},
		encodeRemainLength(ln),
		[]byte{byte(messageID / 256), byte(messageID % 256)},
		qos,
	}, nil)
}

func packUNSUBACK(messageID uint16) []byte {
	return []byte{UNSUBACK * 16, 2, byte(messageID / 256), byte(messageID % 256)}
}

func packPINGRESP() []byte {
	return []byte{PINGRESP * 16, 0}
}

/* Unpack recieve data */

func unpackCONNECT(remaining []byte) (clientID string, willTopic string, willMessage string, loginName string, loginPassword string, err error) {
	connectFlag := remaining[9]
	usernameFlag := (connectFlag & 0x80) != 0
	passwordFlag := (connectFlag & 0x40) != 0
	/*
	   willRetain := (connectFlag & 0x20) != 0
	   qos := byte((connectFlag / 4) & 0x04)
	   willFlag := (connectFlag & 0x02) != 0
	   cleanSession := (connectFlag & 0x01) != 0
	   keepAliveTime := bytes_to_uint16(remaining[10:12])
	*/

	n := 12
	ln := int(bytes_to_uint16(remaining[n : n+2]))
	clientID = bytes_to_str(remaining[n+2 : n+2+ln])

	if usernameFlag {
		n += ln
		ln = int(bytes_to_uint16(remaining[n : n+2]))
		loginName = bytes_to_str(remaining[n+2 : n+2+ln])
	}

	if passwordFlag {
		n += ln
		ln = int(bytes_to_uint16(remaining[n : n+2]))
		loginPassword = bytes_to_str(remaining[n+2 : n+2+ln])
	}

	return
}

func unpackPUBLISH(remaining []byte) (topic string, messageID uint16, payload []byte, err error) {
	ln := bytes_to_uint16(remaining[0:2])
	topic = bytes_to_str(remaining[2 : 2+ln])
	messageID = bytes_to_uint16(remaining[2+ln : 4+ln])
	payload = remaining[4+ln:]
	return
}

func unpackPUBREL(remaining []byte) (messageID uint16, err error) {
	messageID = bytes_to_uint16(remaining[0:2])
	return
}
func unpackPUBACK(remaining []byte) (messageID uint16, err error) {
	messageID = bytes_to_uint16(remaining[0:2])
	return
}

func unpackSUBSCRIBE(remaining []byte) (messageID uint16, topics []string, err error) {
	messageID = bytes_to_uint16(remaining[0:2])
	topics = make([]string, len(remaining)/2)

	n := 0
	for i := 2; i < len(remaining); {
		ln := int(bytes_to_uint16(remaining[i : i+2]))
		topics[n] = bytes_to_str(remaining[i+2 : i+2+ln])
		n++
		i += 2 + ln + 1
	}
	topics = topics[:n]
	return
}

func unpackUNSUBSCRIBE(remaining []byte) (messageID uint16, topics []string, err error) {
	messageID = bytes_to_uint16(remaining[0:2])
	topics = make([]string, len(remaining)/2)

	n := 0
	for i := 2; i < len(remaining); {
		ln := int(bytes_to_uint16(remaining[i : i+2]))
		topics[n] = bytes_to_str(remaining[i+2 : i+2+ln])
		n++
		i += 2 + ln
	}
	topics = topics[:n]

	return
}

func readMessage(conn net.Conn) (command int, fixedHeader []byte, remaining []byte, err error) {
	fixedHeader = make([]byte, 5)
	n, err := conn.Read(fixedHeader[0:1])
	if err != nil && n < 1 {
		err = errors.New("Can't read packets")
	}
	if err != nil {
		return
	}

	var i int
	for i = 0; i < 4; i++ {
		n, err := conn.Read(fixedHeader[i+1 : i+2])
		if err != nil && n < 1 {
			err = errors.New("Can't read packets")
		}
		if err != nil {
			break
		}
		if fixedHeader[i+1] <= 127 {
			break
		}
	}
	if err != nil {
		return
	}

	fixedHeader = fixedHeader[0 : i+2]
	remainLength := decodeRemainLength(fixedHeader[1:])

	remaining = make([]byte, remainLength)
	conn.Read(remaining)

	command = int(fixedHeader[0] / 16)

	return
}
