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

			// Check QA Marker, need Question (0-64)
			if buffer[9] > 64 {
				return
			}

			repeats := int(buffer[8])

			answer := on_accept_qestion(buffer[10:], addr, repeats)
			if answer == nil {
				return
			}

			answ_buffer := make([]byte, len(answer)+10)

			// Put Question ID
			copy(answ_buffer[:8], buffer[:8])

			// Put Repeats
			answ_buffer[8] = buffer[8]

			// Put QA Marker = Answer
			answ_buffer[9] = 255

			// Copy data
			copy(answ_buffer[10:], answer)

			conn.WriteTo(answ_buffer, addr)
		}()
	}
}
