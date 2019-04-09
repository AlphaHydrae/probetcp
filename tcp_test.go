package main

import (
	"net"
	"testing"
	"time"
)

func TestTcp(t *testing.T) {

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()

	go func() {
		_, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}
	}()

	config := &TCPProbeConfig{}
	config.Address = "localhost:8081"
	config.Timeout = time.Duration(1e9)

	result, err := probeTCPEndpoint(config)
	if err != nil {
		t.Fatal(err)
	}

	if !result.Success {
		t.Fatalf("Probe failed %v\n", result)
	}
}
