package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn net.Conn
	msgCh chan Message
	delCh chan *Peer
}

func (p *Peer) Send(msg []byte) (int, error){
	return p.conn.Write(msg)
}
 
func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer {
		conn: conn,
		msgCh: msgCh,
		delCh: delCh,
	}
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delCh <- p
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var cmd Command
		if v.Type() == resp.Array {
			rawCMD := v.Array()[0]
			switch rawCMD.String() {
			case CommandSET:
				if len(v.Array()) != 3 {
					return fmt.Errorf("invalid no of vars for SET command")
				}
				cmd = SetCommand {
					key: v.Array()[1].Bytes(),
					val: v.Array()[2].Bytes(),
				}

			case CommandGET:
				if len(v.Array()) != 2 {
					return fmt.Errorf("invalid no of vars for SET command")
				}
				cmd = GetCommand {
					key: v.Array()[1].Bytes(),
				}

			default:
				fmt.Println("Command not found: ", rawCMD)

			}

			p.msgCh <- Message{
					cmd: cmd,
					peer: p,
				}
		}
	}
	return nil
}