package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	a2 "github.com/adan7950/assignment02IBC"
)

func main() {
	satoshiAddress := ":9000"

	var nodeID string
	fmt.Print("Enter Node ID: ")
	fmt.Scan(&nodeID)
	
	fmt.Print("Assigned Node ID: ", nodeID)

	//Code to connect to Satoshi
	c, err := net.Dial("tcp", satoshiAddress)

	if err != nil {
		log.Fatal(err)
	}

	// Code to send node ID to Satoshi
	c.Write([]byte(nodeID))

	// Code to receive and print the chain
	var chainHead *a2.Block

	dec := gob.NewDecoder(c)
	err = dec.Decode((&chainHead))
	if err != nil {
		log.Println(err)
	}

	a2.ListBlocks(chainHead)
}