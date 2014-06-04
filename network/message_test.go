package network

import "testing"
import "strconv"

func TestValidate(t *testing.T) {
	msg := new(Message)

	// Empty Message
	err := msg.validate()
	if err == nil {
		t.Log("Failed at Empty Message Invalidation.")
		t.Fail()
	}

	// Valid Message
	msg.magic = knownMagic
	msg.payload = []byte{'a', 'b', 'c', 'd'}
	msg.length = 4
	msg.command = "version"
	msg.checksum = 0xd8022f20

	err = msg.validate()
	if err != nil {
		t.Log("Failed at Validating Message: ", err)
		t.Fail()
	}

	// Invalid Magic Number
	msg.magic += 1
	err = msg.validate()
	if err != MessageError(EMAGIC) {
		t.Log("Failed at invalidating magic number: ", err)
		t.Fail()
	}
	msg.magic -= 1

	// Invalid Checksum
	msg.checksum += 1
	err = msg.validate()
	if err != MessageError(ECHECK) {
		t.Log("Failed at invalidating checksum: ", err)
		t.Fail()
	}
	msg.checksum -= 1

	// Invalid Length
	msg.length += 1
	err = msg.validate()
	if err != MessageError(EPALEN) {
		t.Log("Failed at invalidating length: ", err)
		t.Fail()
	}
}

type Payload []byte

func (p Payload) Serialize() []byte {
	return []byte(p)
}

func TestMakeMessage(t *testing.T) {

	str := "version"

	msg := MakeMessage(str, Payload([]byte{'a', 'b', 'c', 'd'}), nil)

	// Basic Validation Test
	err := msg.validate()
	if err != nil {
		t.Log("Created invalid message: ", err)
		t.Fail()
	}

	// Check that values are correct
	str = "notversion"
	if msg.command != strconv.Quote("version") {
		t.Log("Command string doesn't make defensive copy: ", msg.command)
		t.Fail()
	}

	if msg.length != 4 {
		t.Log("Incorrect payload length.")
		t.Fail()
	}

}
