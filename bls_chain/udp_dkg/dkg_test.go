package main

import (
	"testing"
	"net"
)

func TestInitNetwork(t *testing.T) {
	InitNetwork("9999")
	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8964}
	unicast("test", addr)
	broadcast("gogogo")
}

func TestInitBls(t *testing.T) {
	InitBls()
}
