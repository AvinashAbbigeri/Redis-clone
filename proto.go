package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/tidwall/resp"
)

// Represents the possible commands.
const (
	CommandSET    = "set"
	CommandGET    = "get"
	CommandHELLO  = "hello"
	CommandClient = "client"
)

// Serves as a marker for different command types.
type Command interface{}

// Set command with key-value pair and optional TTL.
type SetCommand struct {
	key, val []byte
	ttl      time.Duration // Time-to-live for the key-value pair.
}

// Get command to retrieve value by key.
type GetCommand struct {
	key []byte
}

// Represents client command with string value.
type ClientCommand struct {
	value string
}

// HELLO command with string.
type HelloCommand struct {
	value string
}

// Converts Go map to RESP compliant map response.
func respWriteMap(m map[string]string) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%%%d\r\n", len(m))) // RESP map header
	rw := resp.NewWriter(buf)
	for k, v := range m {
		rw.WriteString(k) // Write key
		rw.WriteString(v) // Write value
	}
	return buf.Bytes()
}
