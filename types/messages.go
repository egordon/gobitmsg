package types

import (
	"time"
	"net"
	"encoding/binary"
)

const (
	netAddrLen = 38
	versionLen = 80
	varIntLen = 9
)

type VarInt uint64

type Hash []byte

type NetAddr struct {
	Time time.Time
	Stream uint32
	Services uint64
	IP net.IP
	Port uint16
}

func (n *NetAddr) Serialize() []byte {
	if n == nil {
		return nil
	}
        ret := make([]byte, 0, netAddrLen)
        binary.BigEndian.PutUint64(uint64(n.Time.Unix()))
        binary.BigEndian.PutUint32(n.Stream)
	binary.BigEndian.PutUint64(n.Services)
	append(ret, []byte(n.IP.To16()))
	binary.BigEndian.PutUint16(n.Port)
	return ret
}

func (n *NetAddr) Unserialize(work []byte) error {
	if len(work) < netAddrLen {
		return SerialError(ELEN)
	}
	if n == nil {
		return SerialEroror(ENIL)
	}

	n.Time = time.Unix(binary.BigEndian.Uint64(work[:8], 0))
	n.Stream = binary.BigEndian.Uint32(work[8:12])
	n.Services = binary.BigEndian.Uint64(work[12:20])
	copy([]byte(n.IP), work[20:36])
	n.Port = binary.BigEndian.Uint16(work[36:38])
	return nil
}

// See https://bitmessage.org/wiki/Protocol_specification#Message_types

type Version struct {
	Version uint32
	Services uint64
	Timestamp time.Time
	Addr_Recv NetAddr
	Addr_From NetAddr
	Nonce uint64
	Streams []uint64
}

func (v *Version) Serialize() []byte {
	ret := make([]byte, 0, versionLen)
	
	binary.BigEndian.PutUint32(ret, v.Version)
	binary.BigEndian.PutUint64(ret, v.Services)
	binary.BigEndian.PutUint64(ret, uint64(v.Timestamp.Unix()))
	net := n.Addr_Recv.Serialize()
	append(ret, net[12:])
	net = n.Addr_From.Serialize()
	append(ret, net[12:])
	binary.BigEndian.PutUint64(ret, v.Nonce)

	return ret
}

type Addr []NetAddr

type Inv []Hash

type GetData Inv
