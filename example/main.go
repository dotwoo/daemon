package main

import (
	"github.com/dotwoo/daemon"
)

func main() {
	ds := NewSample()
	daemon.Run(ds)

}
