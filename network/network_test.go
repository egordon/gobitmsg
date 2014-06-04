package network

import "testing"
import "time"
import "fmt"

func TestNetwork(t *testing.T) {
	// Set up a listening server
	fmt.Println("Setting up listener...")
	err := Listen(":4444")
	if err != nil {
		t.Log("Failed to set up listening server: ", err)
		t.FailNow()
	}

	// Make a peer associated with out server
	fmt.Println("Setting up peer...")
	p := MakePeer("127.0.0.1:4444")
	if p == nil {
		t.Log("Failed to make new peer.")
		t.FailNow()
	}
	
	// Connect to that peer
	fmt.Println("Connecting to Peer...")
	err = ConnectToPeer(p)
	if err != nil {
		t.Log("Could not connect to peer: ")
		t.FailNow()
	}

	// Check that the peer exists and is connected
	fmt.Println("Checking Peer...")
	p2 := GetPeer("127.0.0.1:4444")
	if p2 != p || !p2.IsConnected() {
		t.Log("Peer does not match and is not connected")
		t.FailNow()
	}
	
	// Construct a message for that peer
	fmt.Println("Making Message...")
	m := MakeMessage("version", Payload([]byte{'a', 'b', 'c', 'd'}), p2)
	err = Send(m)
	if err != nil {
		t.Log("Faild to send message: ", err)
		t.FailNow()
	}

	// Check to see that message is received
	fmt.Println("Receiving Message...")
	m2, err := GetNext(2 * time.Second)
	if m2 == nil || err != nil {
		t.Log("Failed to receive message: ", err)
		t.FailNow()
	}

	// Check message integrity
	fmt.Println("Checking Message...")
	err = m2.validate()
	if err != nil {
		t.Log("Message not valid: ", err)
		t.Fail()
	}
	
	test := true
	for i := 0; i < 4; i++ {
		if m.payload[i] != m2.payload[i] {
			test = false
		}
	}

	if !test {
		t.Log("Send message lost integrity: ", m2.payload)
		t.Fail()
	}
}
