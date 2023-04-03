package main

import "goinx/net"

func main() {
	s := net.NewServer("[goinx V0.1]")
	s.Serve()
}
