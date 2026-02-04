package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

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

const (
	CommandSET   = "SET"
	CommandGET   = "GET"
	CommandHELLO = "hello"
)

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("Read %s\n", v.Type())

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSET:
					if len(v.Array()) != 3 {
						return nil, fmt.Errorf("invalid no of vars for SET command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
					return cmd, nil

				case CommandGET:
					if len(v.Array()) != 2 {
						return nil, fmt.Errorf("invalid no of vars for SET command")
					}
					cmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}
					return cmd, nil

				}
			}
		}
		// return nil, fmt.Errorf("invalid or unknown command: %s", raw)
	}
	return nil, fmt.Errorf("invalid or unknown command: %s", raw)
}

// To handle hello command from official redis client
// According to RESP3 Maps Protocol
func writeRespMap(m map[string]string) []byte {
	buf := &bytes.Buffer{}
	rw := resp.NewWriter(buf)
	rw.WriteString("OK")
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("+%s\r\n", k))
		buf.WriteString(fmt.Sprintf(":%s\r\n", v))
	}
	return buf.Bytes()
}
