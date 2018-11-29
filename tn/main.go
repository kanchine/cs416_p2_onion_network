package main

import (
	"fmt"

	"./tornode"
)

func main() {
	fmt.Println("launching tor node")
	tnerr := tornode.InitTorNode("", "", "", ".", ".", 1000)
	fmt.Println(tnerr)
}
