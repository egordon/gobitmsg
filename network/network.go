package network

import (
	"net"
	"log"
)

type Serializer interface {
        Serialize() []byte // Convert object to byte slice
}

const chanBufLen = 100

var peerlist map[string]*Peer
var recvChan chan *Message
var sendChan chan *Message

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

func GetPeer(host string) *Peer {

	checkInit()

	p, ok := peerlist[host]
	if ok {
		return p
	} else {
		return nil
	}
}

func Listen(port string) (chan int, error) {

	checkInit()

	ln, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	quit := make(chan int)

	go accept(ln, quit)

	return quit, nil
}

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

	return p.connect(recvChan)
}
	

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
	


func accept(ln net.Listener, quit chan int) {
	isRunning := true
	for isRunning {
		select {
		case <-quit:
			isRunning = false
		default:
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				conn.Close()
				continue
			}

			peer := MakePeer(conn.RemoteAddr().String())
			peer.conn = conn
			peer.connected = true
			peerlist[peer.hostname] = peer

			go peer.recv(recvChan)
		}
	}
}	

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

