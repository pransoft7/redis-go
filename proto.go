package main

import (
	"bytes"
	"fmt"

	"github.com/tidwall/resp"
)

type Command interface {
	
}

type SetCommand struct {
	key, val []byte
}

type GetCommand struct {
	key []byte
}

type HelloCommand struct {
	value string
}

type ClientCommand struct {
	value string
}

const (
	CommandSET = "set"
	CommandGET = "get"
	CommandHELLO = "hello"
	CommandClient = "client"
)


// To handle hello command from official redis client
// According to RESP3 Maps Protocol
func writeRespMap(m map[string]string) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	rw := resp.NewWriter(buf)
	for k, v := range m {
		rw.WriteString(k)
		rw.WriteString(":" + v)
	}
	return buf.Bytes()
}