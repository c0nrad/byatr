package main

import (
	"bytes"
	"fmt"
	"net"
)

const (
	MaxBlockSize   = 64
	ResponseBuffer = 4048
)

func test() {
	for i := 0; i < MaxBlockSize; i++ {

	}
}

const DebugMode = false

func send(conn net.Conn, in []byte) []byte {
	conn.Write(in)

	buffer := make([]byte, ResponseBuffer)
	count, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	buffer = buffer[:count]
	if DebugMode {
		fmt.Println(in, ":", len(buffer))
	}

	return buffer
}

func preblock(conn net.Conn) int {
	baseLen := len(send(conn, []byte{'C'}))

	for i := 1; i < MaxBlockSize; i++ {
		payload := bytes.Repeat([]byte{'C'}, i)
		response := send(conn, payload)

		if baseLen != len(response) {
			return i
		}

	}
	return -1
}

func blockSize(conn net.Conn) int {
	prevIndex := 0
	prevLen := len(send(conn, []byte{'C'}))
	prevDiff := 0

	for i := 1; i < MaxBlockSize; i++ {
		payload := bytes.Repeat([]byte{'C'}, i)
		response := send(conn, payload)

		if prevLen != len(response) {
			if prevDiff == len(response)-prevLen {
				return i - prevIndex
			}
			prevIndex = i
			prevDiff = len(response) - prevLen
			prevLen = len(response)
		}
	}
	return -1
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	pre := preblock(conn)
	fmt.Println("The pre-block length: ", pre)

	block := blockSize(conn)
	fmt.Println("The block length: ", block)
}
