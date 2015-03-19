package main

import (
	"bytes"
	"fmt"
	"net"
)

const (
	MaxBlockSize   = 64
	ResponseBuffer = 4048
	FillCharacter  = byte('C')
)

var Printable []byte = []byte("abcdefghijklmnopqrstuvwxyzBCDEFGHIJKLMNOPQRSTUVWXYZ ,.!?'1234567890-\n ")

//var Printable []byte = []byte("SECRET")

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

func fillblock(conn net.Conn) int {
	baseLen := len(send(conn, []byte{FillCharacter}))

	for i := 1; i < MaxBlockSize; i++ {
		payload := bytes.Repeat([]byte{FillCharacter}, i)
		response := send(conn, payload)

		if baseLen != len(response) {
			return i
		}

	}
	return -1
}

func blockSize(conn net.Conn) int {
	prevIndex := 0
	prevLen := len(send(conn, []byte{FillCharacter}))
	prevDiff := 0

	for i := 1; i < MaxBlockSize; i++ {
		payload := bytes.Repeat([]byte{FillCharacter}, i)
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

func getBlock(buffer []byte, index, blockSize int) []byte {
	return buffer[index*blockSize : (index+1)*blockSize]
}

func letter(conn net.Conn, prefill, blockSize int, base []byte) byte {

	prefillBuffer := bytes.Repeat([]byte{FillCharacter}, prefill)
	blockBuffer := bytes.Repeat([]byte{FillCharacter}, blockSize-1-len(base))

	guess := append(prefillBuffer, append(blockBuffer, base...)...)
	correct := send(conn, append(prefillBuffer, blockBuffer...))
	block := (len(base) + prefill) / blockSize

	for _, c := range Printable {
		out := append(guess, c)
		guessResult := send(conn, out)

		if bytes.Equal(getBlock(correct, block, blockSize), getBlock(guessResult, block, blockSize)) {
			return c
		}
	}
	return byte(0)
}

func decrypt(conn net.Conn, prefill, blockSize int) []byte {
	base := []byte{}

	for {
		if c := letter(conn, prefill, blockSize, base); c != byte(0) {
			base = append(base, c)
			fmt.Println(string(base))
		} else {
			return base
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	decrypt(conn, 0, 16)
}
