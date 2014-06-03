package network

import (
	"net"
	"io"
	"log"
)

type Peer struct {
	hostname string // Host and Port combination
	conn net.Conn // TCP Connection
	connected bool // Whether connection is alive or not
}

func MakePeer(host string) *Peer {
	peer := new(Peer)
	peer.hostname = host
	peer.connected = false
	peer.conn = nil

	return peer
}

func (p *Peer) Hostname() string {
	if p == nil {
		return ""
	} else {
		return p.hostname
	}
}

func (p *Peer) IsConnected() bool {
	if p.conn == nil {
		return false
	} else {
		return true
	}
}

func (p *Peer) send(msg []byte) error {
	_, err := p.conn.Write(msg)
	return err
}


func (p *Peer) connect(recvChannel chan *Message) (error) {
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
	msgHead := make([]byte, HeaderLen, HeaderLen)

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
