package network

type PeerError int

const (
	ENTFND = iota
	EPRDUP = iota
)

func (errno PeerError) Error() string {
	switch(errno) {
	case ENTFND:
		return "Peer not in Peer List"
	case EPRDUP:
		return "Peer already in list and connected"
	default:
		return "Unkown Peer Error"
	}
}


const (
	ESMALL = iota
	EMAGIC = iota
	EPALEN = iota
	ECHECK = iota
)

type MessageError int

func (errno MessageError) Error() string {
	switch(errno) {
	case ESMALL:
		return "Header size too small"
	case EMAGIC:
		return "Magic number in header is unknown"
	case EPALEN:
		return "Payload length does not match header"
	case ECHECK:
		return "Invalid SHA512 checksum for payload"
	default:
		return "Unknown Message Error"
	}
}
