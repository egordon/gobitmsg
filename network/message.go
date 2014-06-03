package network

import (
	"encoding/binary"
	"strconv"
	"crypto/sha512"
)

const (
	HeaderLen = 24
	knownMagic = 0xe9beb4d9
)

type Message struct {
	peer *Peer // Peer who sent the message (nil if local is sending)
	magic uint32 // Magic number associated with network
	command string // Action this message wants to take
	length uint32 // Length of the payload
	checksum uint32 // First 4 bytes of sha512(payload)
	payload []byte
}

func MakeMessage(cmd string, pload Serializer, recipient *Peer) *Message {
	msg := new(Message)
	msg.peer = recipient
	msg.magic = knownMagic
	msg.payload = pload.Serialize()
	msg.length = uint32(len(msg.payload))
	msg.command = strconv.QuoteToASCII(cmd)
	
	digest := sha512.Sum512(msg.payload)
	checksum := make([]byte, 4, 4)
	for i := 0; i < 4; i++ {
		checksum[i] = digest[i]
	}
	msg.checksum = binary.BigEndian.Uint32(checksum)
	return msg
}


func (m *Message) Serialize() []byte {
	ret := make([]byte, 0, HeaderLen)
	binary.BigEndian.PutUint32(ret, m.magic)
	ret = strconv.AppendQuoteToASCII(ret, strconv.QuoteToASCII(m.command))

	// Ensure string had a max size of 12
	ret = ret[:16]
	binary.BigEndian.PutUint32(ret, m.length)
	binary.BigEndian.PutUint32(ret, m.checksum)
	ret = append(ret, m.payload...)

	return ret
}

func (m *Message) validate() error {
	
	if m.magic != knownMagic {
		return MessageError(EMAGIC)
	}
	
	if m.length != uint32(len(m.payload)) {
		return MessageError(EPALEN)
	}

	checksum := make([]byte, 4, 4)
	digest := sha512.Sum512(m.payload)

	checksum[0] = digest[0]
	checksum[1] = digest[1]
	checksum[2] = digest[2]
	checksum[3] = digest[3]

	if m.checksum != binary.BigEndian.Uint32(checksum) {
		return MessageError(ECHECK)
	}

	return nil
}

func (m *Message) makeHeader(rawBytes []byte) error {
	if len(rawBytes) < HeaderLen {
		return MessageError(ESMALL)
	}
	m.magic = binary.BigEndian.Uint32(rawBytes[:4])
	m.command = string(rawBytes[4:16])
	m.length = binary.BigEndian.Uint32(rawBytes[16:20])
	m.checksum = binary.BigEndian.Uint32(rawBytes[20:24])
	m.payload = nil
	return nil
}

func (m *Message) setPayload(rawBytes []byte) {
	m.payload = make([]byte, len(rawBytes), len(rawBytes))
	copy(m.payload, rawBytes)
}

func (m *Message) Payload() []byte {
	if m == nil {
		return nil
	} else {
		return m.payload
	}
}

func (m *Message) Command() string {
	if m == nil {
		return ""
	} else {
		return m.command
	}
}

func (m *Message) Sender() *Peer {
	if m == nil {
		return nil
	} else {
		return m.peer
	}
}
