package types

import "fmt"

type VarIntError int

const (
	EINVLD = iota
)

func (e VerIntError) Error() string {
	switch e {
	case EINVLD:
		return "Invalid Variable Integer Encoding"
	default:
		return fmt.Sprintln("VarInt has length: ", int(e)) 
	}
}
