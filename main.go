package main

import (
	"fmt"
	"net"
	"strconv"
)

func main() {
	go func() {
		fmt.Println(Start_Server("localhost:1083",
			func(question []byte, addr net.Addr, repeats int) (answer []byte) {
				fmt.Println(string(question) + " / repeats: " + strconv.Itoa(repeats))
				return []byte("Hello!")
			}))
	}()

	fmt.Println("Send datagram")
	answ, err, rep := Send_question("localhost:1083", []byte("Hello World."))

	fmt.Println(err)
	fmt.Println(answ)
	fmt.Println(rep)
}
