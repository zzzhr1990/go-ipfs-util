package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zzzhr1990/go-ipfs-util/adder"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hashcalc <hash>")
		os.Exit(1)
	}
	filePath := os.Args[1]
	fileNode, err := adder.HashFile(context.Background(), "", filePath)
	if err != nil {
		print("err!")
		println(err.Error())
		return
	}
	for _, link := range fileNode.Links() {
		println(">")
		fmt.Println(link.Name)
		fmt.Println(link.Cid.String())
	}
	println("----------------------------------------------------")

	fmt.Println(fileNode.Cid().String())
	//cs, _ := hashcalc.CalcFileHash(filePath, 0)
	//println("----------------------------------------------------")
	//println(cs)
}
