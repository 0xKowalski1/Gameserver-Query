package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

var serverAddr string = "join.insanitycraft.net"
var serverPort int = 25565

func main() {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", serverAddr, serverPort), 5*time.Second)
	if err != nil {
		panic(err)
	}

	// Handshake
	data := []byte{0x00}
	data = append(data, encodeVarInt(766)...) // Protocol version, 766 corresponds to Minecraft 1.20.5
	data = append(data, encodeString(serverAddr)...)
	data = append(data, encodeVarInt(serverPort)...) // Default Minecraft port
	data = append(data, encodeVarInt(1)...)          // Next state: status

	// Send the data with length prefix
	conn.Write(encodeVarInt(len(data)))
	conn.Write(data)

	// Request Status
	data = []byte{0x00}
	conn.Write(encodeVarInt(len(data)))
	conn.Write(data)

	_, err = readVarInt(conn)
	if err != nil {
		fmt.Printf("Error reading response length: %v\n", err)
		return
	}

	packetID, err := readVarInt(conn)
	if err != nil {
		fmt.Printf("Error reading packet ID: %v\n", err)
		return
	}

	if packetID != 0 { // We expect the packet ID for a response to be 0
		fmt.Printf("Unexpected packet ID %d\n", packetID)
		return
	}

	jsonData, err := readString(conn)
	if err != nil {
		fmt.Printf("Error reading JSON data: %v\n", err)
		return
	}

	fmt.Println("Server response:", jsonData)
}

func encodeVarInt(value int) []byte {
	var buf []byte
	for {
		part := byte(value & 0x7F)
		value >>= 7
		if value != 0 {
			part |= 0x80
		}
		buf = append(buf, part)
		if value == 0 {
			break
		}
	}
	return buf
}

func encodeString(s string) []byte {
	return append(encodeVarInt(len(s)), s...)
}

func readString(conn net.Conn) (string, error) {
	length, err := readVarInt(conn)
	if err != nil {
		return "", err
	}
	bytes := make([]byte, length)
	_, err = bufio.NewReader(conn).Read(bytes)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func readVarInt(conn net.Conn) (int, error) {
	var value int
	var size int
	var b byte
	var err error

	for {
		b, err = readByte(conn)
		if err != nil {
			return 0, err
		}
		value |= int(b&0x7F) << (7 * size)
		size++
		if size > 5 {
			return 0, fmt.Errorf("VarInt is too big")
		}
		if (b & 0x80) == 0 {
			break
		}
	}
	return value, nil
}
func readByte(conn net.Conn) (byte, error) {
	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}
