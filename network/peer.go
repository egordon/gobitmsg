package network

import (
	"io"
	"log"
	"net"
)

// Peer represents and holds connection information for an external connection specified by a 
// hostname/port combination.
type Peer struct {
	hostname  string   // Host and Port combination
	conn      net.Conn // TCP Connection
	connected bool     // Whether connection is alive or not
}

// MakePeer creates an unconnected peer from a host/port combination
func MakePeer(host string) *Peer {
	peer := new(Peer)
	peer.hostname = host
	peer.connected = false
	peer.conn = nil

	return peer
}

// Hostname returns the host/port combination for a given peer
func (p *Peer) Hostname() string {
	if p == nil {
		return ""
	} else {
		return p.hostname
	}
}

// IsConnected returns the connection status of the peer
func (p *Peer) IsConnected() bool {
	if p.conn == nil {
		return false
	} else {
		return true
	}
}

// send (unexported) sends a raw byte slice to the peer
func (p *Peer) send(msg []byte) error {
	_, err := p.conn.Write(msg)
	return err
}

func (p *Peer) connect(recvChannel chan *Message) error {
	if p.IsConnected() {
		return nil
	}

	conn, err := net.Dial("tcp", p.hostname)
	if err != nil {
		conn.Close()
		conn = nil
		return err
	}

	p.conn = conn
	p.connected = true

	// Set up receive routine
	go p.recv(recvChannel)
	return nil
}

func (p *Peer) recv(c chan *Message) {
	msgHead := make([]byte, headerLen, headerLen)

	for p.IsConnected() {
		msg := new(Message)

		_, err := io.ReadFull(p.conn, msgHead)
		if err != nil {
			break
		}

		msg.makeHeader(msgHead)

		msgPayload := make([]byte, msg.length, msg.length)

		_, err = io.ReadFull(p.conn, msgPayload)
		if err != nil {
			break
		}

		msg.setPayload(msgPayload)

		err = msg.validate()
		if err != nil {
			log.Println("Peer %s sent invalid Message: %s", p, err)
		}

		msg.peer = p

		// Pass valid message along for handling
		c <- msg
	}
	msgHead = nil
}
