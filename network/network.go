package network

type Serializer interface {
        Serialize() []byte // Convert object to byte slice
}
