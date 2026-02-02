package main

import (
	"github.com/brezzgg/delease/cmd"
	"github.com/brezzgg/go-packages/lg"
)

func main() {
	defer lg.Close()

	cmd.Run()
}
