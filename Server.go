package main

import (
	"fmt"
	"net"
)

func Start_Server(listen_addr string, on_accept_qestion func(question []byte, addr net.Addr, repeats int) (answer []byte)) error {
	conn, err := net.ListenPacket("udp", listen_addr)
	if err != nil {
		return err
	}

	buffer := make([]byte, 65535)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Ошибка при чтении данных:", err)
			continue
		}

		go func() {
			buffer = buffer[:n]

			// Check Question Marker
			if buffer[0] != 94 || buffer[1] != 48 {
				return
			}

			// Check Adaptive Q Marker
			if buffer[2] != buffer[10]*4 {
				fmt.Println("AQ Marker invalid.")
				return
			}

			repeats := int(buffer[11])

			answer := on_accept_qestion(buffer[12:], addr, repeats)
			if answer == nil {
				return
			}

			answ_buffer := make([]byte, len(answer)+12)

			// Answer Marker
			answ_buffer[0] = 92
			answ_buffer[1] = 46

			// Adapted Answ Marker
			answ_buffer[2] = buffer[10] * 4

			// Put repeats value
			answ_buffer[11] = buffer[11]

			// Question ID value
			copy(answ_buffer[3:], buffer[3:11])

			// Copy data
			copy(answ_buffer[12:], answer)

			conn.WriteTo(answ_buffer, addr)
		}()
	}
}
