package types

// Serializer denotes a type that can be converted to a byte slice
// to be sent over the network.                                                                                                                                     
type Serializer interface {
        Serialize() []byte // Convert object to byte slice 
}

type Unserializer interface {
	Unserialize([]byte) // Fill object using byte slice
}

type FullSerializer interface {
	Serializer
	Unserializer
}
