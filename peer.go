package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

// Represented connected client.
type Peer struct {
	conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
}

// Sends the messaage to to client and
// returns number of bytes written.
func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

// Creates a new peer instance.
func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
		delCh: delCh,
	}
}

// Reads the data from the client if RESP format.
func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delCh <- p
			break
		}
		if err != nil {
			log.Println("Read error:", err)
			log.Fatal(err)
		}

		// RESP protocol represtents commands as array.
		var cmd Command
		if v.Type() == resp.Array {
			rawCMD := v.Array()[0]
			switch rawCMD.String() {

			// Used for client related commands.
			case CommandClient:
				cmd = ClientCommand{
					value: v.Array()[1].String(),
				}

			// Retrieves the key from the key-values store.
			case CommandGET:
				cmd = GetCommand{
					key: v.Array()[1].Bytes(),
				}

			// Stores a key-value pair.
			case CommandSET:
				cmd = SetCommand{
					key: v.Array()[1].Bytes(),
					val: v.Array()[2].Bytes(),
				}

			// Resonse to HELLO command.
			case CommandHELLO:
				cmd = HelloCommand{
					value: v.Array()[1].String(),
				}

			// For unrecognized commands.
			default:
				fmt.Println("got this unhandled command", rawCMD)
			}

			// Sends the parse command and
			// the peer refrence to the server's message channel.
			p.msgCh <- Message{
				cmd:  cmd,
				peer: p,
			}
		}
	}
	return nil
}
