package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	a2 "github.com/adan7950/assignment02IBC"
)


type Node struct {
	NodeID string
	conn net.Conn
}

var clientArray = make([]Node , 0 , 20)
var chainHead *a2.Block


func centralRoutine(clientCh chan Node){
	fmt.Println("Started Server central routine")
	
	for {
		newClient := <-clientCh

		// Append to client array
		clientArray = append(clientArray, newClient)

		// Giving some coins to the reader
		newBlock := []a2.BlockData{{Title: "Satoshi2" + newClient.NodeID, Sender: "Satoshi", Receiver: newClient.NodeID, Amount: 10}}
		chainHead = a2.InsertBlock(newBlock, chainHead)
	
		a2.ListBlocks(chainHead)
		
		// Sending to all clients
		for it:=0 ; it<len(clientArray) ; it++{
			// Making a bufio encoder
			gobEncoder := gob.NewEncoder(clientArray[it].conn)
			err := gobEncoder.Encode(chainHead)
			if err != nil {
				// log.Println(err)
			}	
		}
	}
}

func satoshiClientHandler( conn net.Conn, clientCh chan Node ){
	buffer := make([]byte, 1000)
	n,_ := conn.Read(buffer)
	
	nodeID := string(buffer[0:n])
	
	fmt.Println("Client Received with ID: ", nodeID)
	fmt.Println()
	
	var newClient Node
	newClient.conn = conn
	newClient.NodeID = nodeID
	
	clientCh <- newClient
}


func main() {
	satoshiAddress := ":9000"
	var clientChannel = make(chan Node)
	ln,err := net.Listen("tcp", satoshiAddress)
	if err != nil {
		log.Fatal(err)
	}
	
	//Premining for satoshi
	chainHead = a2.PremineChain(chainHead, 1)

	go centralRoutine(clientChannel)

	fmt.Println("Server Started")
	fmt.Println()

	for {
		conn,err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go satoshiClientHandler(conn, clientChannel)
	}
}