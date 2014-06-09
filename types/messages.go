package types

import (
	"time"
	"net"
	"encoding/binary"
)

const (
	netAddrLen = 26
	versionLen = 80
	varIntLen = 9
)

type VarInt uint64

func (v *VarInt) Serialize() []byte {
	work := make([]byte, varIntLen, varIntLen)
	work[0] = 0xff
	binary.BigEndian.PutUint64(work[1:], uint64(*v))
	if (work[1] | work[2] | work[3] | work[4]) == 0 { // Uint32
		work = work[4:]
		work[0] = 0xfe
		if (work[1] | work[2]) == 0 { // Uint16
			work = work[2:]
			work[0] = 0xfd
			if work[1] == 0 && work[0] < 0xfd { // Uint8
				work = work[0:1]
			}
		}
	}

	ret := make([]byte, len(work), len(work))
	copy(ret, work)
	return ret
}

type Hash []byte

type NetAddr struct {
	Time time.Time
	Stream uint32
	Services uint64
	IP net.IP
	Port uint16
}

// See https://bitmessage.org/wiki/Protocol_specification#Message_types

type Version struct {
	Version int32
	Services uint64
	Timestamp time.Time
	Addr_Recv NetAddr
	Addr_From NetAddr
	Nonce uint64
	Streams []uint64
}

func (v *Version) Serialize() []byte {
	ret := make([]byte, versionLen, versionLen)
	return ret
}

type Addr []NetAddr

type Inv []Hash

type GetData Inv
