package db

import (
	"github.com/egordon/gobitmsg/network"
	"github.com/mxk/go-sqlite/sqlite3"
	"log"
)

var dbConn sqlite3.Conn

func ConnectToDB(fileName string) error {
	
}

func SetInventory(hash []byte, payload network.Serializer) error {
}

func GetInventory(hash []byte) (network.Serializer, error) {
}
