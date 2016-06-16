package main

import (
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/hse-chat/hse-chat-api/HseMsg"
)

// ProtoConnection is a chunked connection
type ProtoConnection struct {
	conn ChunkedConnection
}

func (pConn *ProtoConnection) Read(msg *HseMsg.Request) error {
	chunk, err := pConn.conn.Read()
	if err != nil {
		return err
	}

	err = proto.Unmarshal(chunk, msg)
	if err != nil {
		return err
	}

	log.Printf("Received %s", msg)
	return nil
}

func (pConn *ProtoConnection) Write(msg *HseMsg.ServerMessage) error {
	log.Printf("Sending %s", msg)

	chunk, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}

	return pConn.conn.Write(chunk)
}

// Close closes connection
func (pConn *ProtoConnection) Close() {
	log.Print("Closing connection")
	pConn.conn.Close()
}

// ProtoListener is a type representing chunked listener
type ProtoListener struct {
	listener ChunkedListener
}

func (ln *ProtoListener) listen() error {
	return ln.listener.listen()
}

// Accept accepts new connection and return it
func (ln *ProtoListener) Accept() (ProtoConnection, error) {
	conn, err := ln.listener.Accept()
	if err != nil {
		return ProtoConnection{ChunkedConnection{nil}}, err
	}

	log.Print("Accepted connection")
	chConn := ProtoConnection{conn}

	return chConn, nil
}

// ProtoListen creates chunked listener
func ProtoListen(netType string, addr string) (*ProtoListener, error) {
	chListen, err := ChunkedListen(netType, addr)
	if err != nil {
		return nil, err
	}

	return &ProtoListener{chListen}, nil
}
