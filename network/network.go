// Package network provides threaded P2P functionality for sending
// Messages wrapped in the bitmsg header.
package network

import (
	"log"
	"net"
	"time"
	"fmt"
)

// Default receive/send channel buffer length
const chanBufLen = 100

// Known list of peers.
var peerlist map[string]*Peer

// Channels acting as queues for sending and
// receiving messages.
var recvChan chan *Message
var sendChan chan *Message

// checkInit (unExported) just makes sure that all the global
// scope variables are initialized.
func checkInit() {
	if peerlist == nil {
		peerlist = make(map[string]*Peer)
	}
	if recvChan == nil {
		recvChan = make(chan *Message, chanBufLen)
	}
	if sendChan == nil {
		sendChan = make(chan *Message, chanBufLen)
		go sendThread()
	}

}

// GetPeer checks the peer list to see if a peer exists
// under a hostname. It then returns that peer.
func GetPeer(host string) *Peer {

	checkInit()

	p, ok := peerlist[host]
	if ok {
		return p
	} else {
		return nil
	}
}

// Listen sets up a local TCP BitMSG server on the provided
// port. The port must be in the format ":###", with the colon.
func Listen(port string) (error) {

	checkInit()

	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	go accept(ln)

	return nil
}


// ConnectToPeer adds the provided peer to the peerlist,
// or just connects to it if it already exists. Duplicate
// connected peers are not allowed.
func ConnectToPeer(p *Peer) error {
	_, ok := peerlist[p.hostname]
	if !ok {
		peerlist[p.hostname] = p
	} else {
		p = peerlist[p.hostname]
	}

	if p.IsConnected() {
		return PeerError(EPRDUP)
	}

	err := p.connect(recvChan)
	if err != nil {
		delete(peerlist, p.hostname)
	}
	return err
}


// Send adds a validated messge to the send channel.
func Send(msg *Message) error {

	checkInit()

	err := msg.validate()
	if err != nil {
		return err
	}
	_, ok := peerlist[msg.peer.hostname]
	if ok {
		sendChan <- msg
	} else {
		log.Println("Error, peer not in peerlist")
		return PeerError(ENTFND)
	}
	return nil
}

// GetNext returns the next message in the queue
// or nil if the timeout is reached
func GetNext(timeout time.Duration) (*Message, error) {
	t := make(chan bool, 1)
	
	go func() {
		time.Sleep(timeout)
		t <- true
	}()
	
	select {
	case msg := <-recvChan:
		return msg, nil
	case <-t:
		return nil, MessageError(ETIMEO)
	}
}

// accept (unexported), run as a goroutine, provides the
// listen loop for new peers.
func accept(ln net.Listener) {

	for ln != nil {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("EEEEEEEERRRRRRRRRROOOOOOOOORRRRRRR!")
			log.Println(err)
			conn.Close()
			break
		}
		
		peer := MakePeer(conn.RemoteAddr().String())
		peer.conn = conn
		peer.connected = true
		peerlist[peer.hostname] = peer
		
		go peer.recv(recvChan)
	}
}

// sendThread (unexported), meant to be run as a goroutine,
// continuously pulls off the send Channel, ensures that the
// recipient is still connected, and sends the message.
func sendThread() {
	for {
		msg := <-sendChan
		p, ok := peerlist[msg.peer.hostname]
		if ok {
			p.send(msg.Serialize())
		} else {
			log.Println("Error in sendThread, peer no found")
		}
	}
}
