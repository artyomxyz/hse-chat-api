package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

// ChunkedConnection is a chunked connection
type ChunkedConnection struct {
	conn net.Conn
}

type connectionEndError struct{}

func (e *connectionEndError) Error() string {
	return "Connection ended"
}

func (chConn *ChunkedConnection) Read() ([]byte, error) {
	reader := bufio.NewReader(chConn.conn)
	lenBuffer := make([]byte, 4)
	bytesRead, err := io.ReadFull(reader, lenBuffer)
	if err != nil {
		return nil, err
	}
	if bytesRead != 4 {
		return nil, &connectionEndError{}
	}

	len := binary.LittleEndian.Uint32(lenBuffer)

	messageBuff := make([]byte, len)
	bytesRead, err = io.ReadFull(reader, messageBuff)
	if err != nil {
		return nil, err
	}
	if bytesRead != int(len) {
		return nil, &connectionEndError{}
	}

	return messageBuff, nil
}

func (chConn *ChunkedConnection) Write(chunk []byte) error {
	lenBuffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBuffer, uint32(len(chunk)))

	bytesWritten, err := chConn.conn.Write(lenBuffer)
	if err != nil {
		return err
	}
	if bytesWritten != 4 {
		return &connectionEndError{}
	}

	bytesWritten, err = chConn.conn.Write(chunk)
	if err != nil {
		return err
	}
	if bytesWritten != len(chunk) {
		return &connectionEndError{}
	}

	return nil
}

// Close closes connection
func (chConn *ChunkedConnection) Close() {
	chConn.conn.Close()
}

// ChunkedListener is a type representing chunked listener
type ChunkedListener struct {
	netType  string
	addr     string
	listener net.Listener
}

func (ln *ChunkedListener) listen() error {
	listener, err := net.Listen(ln.netType, ln.addr)

	if err != nil {
		return err
	}

	ln.listener = listener
	return nil
}

// Accept accepts new connection and return it
func (ln *ChunkedListener) Accept() (ChunkedConnection, error) {
	conn, err := ln.listener.Accept()
	if err != nil {
		return ChunkedConnection{nil}, err
	}

	chConn := ChunkedConnection{conn}

	return chConn, nil
}

// ChunkedListen creates chunked listener
func ChunkedListen(netType string, addr string) (ChunkedListener, error) {
	ln := ChunkedListener{netType, addr, nil}

	if err := ln.listen(); err != nil {
		return ln, err
	}

	return ln, nil
}
