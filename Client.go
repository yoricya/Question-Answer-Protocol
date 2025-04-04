package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/rand/v2"
	"net"
	"time"
)

func Send_question(address string, data []byte) (answer []byte, err error, repeats int) {
	if len(data) > 65100 {
		return nil, errors.New("datagram packet length overflow"), 0
	}

	question_id := make([]byte, 8)
	binary.BigEndian.PutUint64(question_id, rand.Uint64())

	datagram := make([]byte, len(data)+12)

	// Question Marker
	datagram[0] = 94
	datagram[1] = 48

	// Adapted Q Marker
	datagram[2] = question_id[7] * 4

	// Question ID value
	copy(datagram[3:], question_id)

	// Copy data
	copy(datagram[12:], data)

	for i := 0; i < 12; i++ {
		if i != 0 {
			time.Sleep(time.Millisecond * 250)
		}

		// Repeats Value
		datagram[11] = byte(i)

		udpAddr, err := net.ResolveUDPAddr("udp", address)
		if err != nil {
			return nil, err, 0
		}

		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			return nil, err, 0
		}

		conn.SetReadDeadline(time.Now().Add(4500 * time.Millisecond))

		// Send datagram
		_, err = conn.Write(datagram)
		if err != nil {
			conn.Close()
			continue
		}

		buffer := make([]byte, 65535)
		for io := 0; io < 10; io++ {
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				conn.Close()
				break
			}
			buffer = buffer[:n]

			// Check Answer Marker
			if buffer[0] != 92 || buffer[1] != 46 {
				conn.Close()
				break
			}

			// Check Adaptive Answer Marker
			if buffer[2] != datagram[2] {
				continue
			}

			// Check Question ID
			if !bytes.Equal(buffer[3:11], datagram[3:11]) {
				continue
			}

			return buffer[12:], nil, int(buffer[11])
		}

		conn.Close()
	}

	return nil, errors.New("server timed out"), 0
}
