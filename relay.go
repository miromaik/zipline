package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

const relayAddr = "localhost:8888"

type Message struct {
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
	Data    []byte `json:"data,omitempty"`
	Done    bool   `json:"done,omitempty"`
	Confirm bool   `json:"confirm,omitempty"`
	Info    string `json:"info,omitempty"`
	Error   string `json:"error,omitempty"`
}

type RelayServer struct {
	waiting map[string]net.Conn
	mu      sync.Mutex
}

func NewRelayServer() *RelayServer {
	return &RelayServer{
		waiting: make(map[string]net.Conn),
	}
}

func (r *RelayServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Relay server listening on %s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v\n", err)
			continue
		}
		go r.handleConnection(conn)
	}
}

func (r *RelayServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	var msg Message
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("Decode error: %v\n", err)
		return
	}

	r.mu.Lock()
	peer, exists := r.waiting[msg.Code]
	if !exists {
		r.waiting[msg.Code] = conn
		r.mu.Unlock()

		io.Copy(io.Discard, conn)

		r.mu.Lock()
		if r.waiting[msg.Code] == conn {
			delete(r.waiting, msg.Code)
		}
		r.mu.Unlock()
		return
	}
	delete(r.waiting, msg.Code)
	r.mu.Unlock()

	go relay(conn, peer)
	relay(peer, conn)
}

func relay(dst, src net.Conn) {
	io.Copy(dst, src)
	dst.Close()
	src.Close()
}

func connectToRelay(code string) (net.Conn, error) {
	conn, err := net.Dial("tcp", relayAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to relay: %w", err)
	}

	msg := Message{Type: "connect", Code: code}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&msg); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send code: %w", err)
	}

	return conn, nil
}

func sendMessage(conn net.Conn, msg Message) error {
	encoder := json.NewEncoder(conn)
	return encoder.Encode(&msg)
}

func receiveMessage(conn net.Conn) (*Message, error) {
	var msg Message
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&msg)
	return &msg, err
}
