package main

import "goinx/gnet"

func main() {
	s := gnet.NewServer("[goinx V0.1]")
	s.Serve()
}
