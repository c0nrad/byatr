package main

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"net"
)

const (
	BlockSize = 16
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}

	fmt.Println("Server listening on :8080")
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	secret := []byte("SECRET")
	key := []byte("ENCRYPTIONKEY123")

	for {
		buffer := make([]byte, 4048)
		count, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Error", err)
			return
		}

		buffer = bytes.TrimSpace(buffer[:count])

		buffer = append(buffer, secret...)
		buffer = pad(buffer)

		fmt.Println("Encrypting", buffer)

		ciphertext := encrypt(buffer, key)
		encoded := encode(ciphertext)
		conn.Write(encoded)
	}
}

func pad(in []byte) []byte {
	extra := BlockSize - (len(in) % BlockSize)
	out := append(in, bytes.Repeat([]byte{0}, extra)...)
	return out
}

func encode(in []byte) []byte {
	size := hex.EncodedLen(len(in))

	out := make([]byte, size)

	hex.Encode(out, in)
	return out
}

func encrypt(buffer, key []byte) []byte {
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(buffer)%16 != 0 {
		panic("Invalid buffer length")
	}

	out := make([]byte, len(buffer))
	for x := 0; x < len(buffer)/16; x += 1 {
		c.Encrypt(out[x*16:(x+1)*16], buffer[x*16:(x+1)*16])
	}

	return out
}
