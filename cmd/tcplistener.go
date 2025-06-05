package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("connection open")
		ch := getLinesChannel(connection)
		for s := range ch {
			fmt.Printf("read: %s\n", s)
		}
		fmt.Println("connection closed")
	}
	// file, err := os.Open("messages.txt")

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// ch := getLinesChannel(file)

	// for s := range ch {
	// 	fmt.Printf("read: %s\n", s)
	// }

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		var build string

		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")

			for i := range len(parts) - 1 {
				ch <- build + parts[i]
				build = ""
			}
			build += parts[len(parts)-1]
		}

		// Emit the final build if any content remains
		if build != "" {
			ch <- build
		}
	}()

	return ch
}
