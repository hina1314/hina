package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const Prompt = "hina> "

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	go readLoop(reader)
	go writeLoop(conn)
	select {}
}

func readLoop(reader *bufio.Reader) {
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}

		fmt.Print(response)
		fmt.Print(Prompt)
	}
}

func writeLoop(conn net.Conn) {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(Prompt)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		_, err = conn.Write([]byte(input + "\n")) // 添加换行符确保服务器能识别命令结束
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}
	}
}
