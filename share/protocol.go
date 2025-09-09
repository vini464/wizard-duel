package share

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// this module standardizes the communication between client and server

type Message struct {
	Type      string `json:"type"` // it can be a command or a status like (login or OK)
	Data      []byte `json:"data,omitempty"`
	Uuid      string `json:"uuid,omitempty"`
	TimeStamp int64  `json:"timestamp"`
}

// message constants
const (
	LOGIN      = "LOGIN"
	LOGOUT     = "LOGOUT"
	REGISTER   = "REGISTER"
	GETBOOSTER = "GETBOOSTER"
	SAVEDECK   = "SAVEDECK"
	PLAY       = "PLAY"
	PLACECARD  = "PLACECARD"
	SKIPPHASE  = "SKIPPHASE"
	SURRENDER  = "SURRENDER"
	OK         = "OK"
	ERROR      = "ERROR"
	INQUEUE    = "INQUEUE"
	PLAYING    = "PLAYING"
)

// communication constants
const (
	SERVERTYPE = "tcp"
	SERVERNAME = "localhost"
	SERVERPORT = "8080"
)

// this function sends a message through a connection and return any occourred error
func SendMessage(conn net.Conn, message Message) error {
	message.TimeStamp = time.Now().UnixNano()
	serialized, err := json.Marshal(message)
	if err != nil {
		return err
	}
	header := make([]byte, 4)
	message_size := len(serialized)
	binary.BigEndian.PutUint32(header, uint32(message_size))

	_, err = conn.Write(header)
	if err != nil {
		return err
	}

	_, err = conn.Write(serialized)

	return err
}

// this function recieves a message through connection saves the info in the given pointer and return any occourred error
func ReceiveMessage(conn net.Conn, message *Message) error {
	header := make([]byte, 4)
	bytes_received := 0
	for bytes_received < len(header) {
		readed, err := conn.Read(header[bytes_received:])
		if err != nil {
			return err
		}
		bytes_received += readed
	}
	data_length := int(binary.BigEndian.Uint32(header))
	data := make([]byte, data_length)
	bytes_received = 0
	for bytes_received < len(data) {
		readed, err := conn.Read(data[bytes_received:])
		if err != nil {
			return err
		}
		bytes_received += readed
	}
	currentTime := time.Now()
	err := json.Unmarshal(data, message)
	if err == nil {
		message.TimeStamp = currentTime.UnixNano() - message.TimeStamp
		fmt.Println("[LATENCY] -", message.TimeStamp, "ns")
	}
	return err
}
