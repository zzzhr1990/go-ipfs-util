package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zzzhr1990/go-ipfs-util/datastores"
	"github.com/zzzhr1990/go-ipfs-util/simple"
)

func Xmain() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hashcalc <hash>")
		os.Exit(1)
	}
	path := os.Args[1]
	println(path)
	info, err := os.Stat(path)
	println("file size: ", info.Size())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//fmt.Println(path, info.Size())
	if !info.IsDir() {
		nd, is, err := simple.Add(path)
		if err != nil {
			println(err.Error())
			os.Exit(2)
		}
		xlHash1 := hex.EncodeToString(is.SumXL())
		xlHash2, err := datastores.CalcFileGcid(path)
		if err != nil {
			println(err.Error())
			os.Exit(3)
		}
		if xlHash1 != xlHash2 {
			println("not equal")
			println(xlHash1)
			println(xlHash2)
			println(path)
			os.Exit(3)
		}

		etag, _ := datastores.GetEtag(path)
		if is.SumWcsEtag() != etag {
			println("not equal etag")
			println(etag)
			println(is.SumWcsEtag())
			println(path)
			os.Exit(3)
		}
		var ct int64 = 0
		var cs int = 0
		is.Reslove(nd.Cid().String(), &ct, &cs)

	}
	//fileNode, err := adder.HashFile(context.Background(), "", filePath)
	//if err != nil {
	//	print("err!")
	//	println(err.Error())
	//	return
	//}
	// fileNode.Links()
	//for _, link := range fileNode.Links() {
	//	println(">")
	//	fmt.Println(link.Name)
	//	fmt.Println(link.Cid.String())
	// fmt.Println(link.c)
	// linkNode, err := fileNode.(context.Background(), link.Cid)
	//}
	//println("----------------------------------------------------")

	//fmt.Println(fileNode.Cid().String())
	//cs, _ := hashcalc.CalcFileHash(filePath, 0)
	//nd, r, err := simple.Add(filePath)
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	//fmt.Println(nd.Cid().String())

	//var ct int64 = 0
	//var cs int = 0
	//println("----------------------------------------------------")
	//r.Reslove(nd.Cid().String(), &ct, &cs)

	//println(ct, " fs: ", r.GetFileSize(), " block size: ", r.GetBlockCount())
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hashcalc <hash>")
		os.Exit(1)
	}
	filePath := os.Args[1]
	fi, err := os.Stat(filePath)
	if err != nil {
		println(err.Error())
		return
	}
	if !fi.IsDir() {
		println("not dir")
		return
	}

	err = filepath.Walk(filePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				println(path, " has errors: ", err.Error())
				return nil
			}
			//fmt.Println(path, info.Size())
			if !info.IsDir() && info.Mode()&os.ModeSymlink != os.ModeSymlink {
				nd, is, err := simple.Add(path)
				if err != nil {
					println(err.Error())
					println(info.IsDir())
					return nil
				}
				xlHash1 := hex.EncodeToString(is.SumXL())
				xlHash2, err := datastores.CalcFileGcid(path)
				if err != nil {
					println(err.Error())
					return err
				}
				if xlHash1 != xlHash2 {
					println("not equal")
					println(xlHash1)
					println(xlHash2)
					println(path)
					os.Exit(3)
				}

				etag, _ := datastores.GetEtag(path)
				if is.SumWcsEtag() != etag {
					println("not equal etag")
					println(etag)
					println(is.SumWcsEtag())
					println(path)
					os.Exit(3)
				}
				var ct int64 = 0
				var cs int = 0
				is.Reslove(nd.Cid().String(), &ct, &cs)

			}
			return nil
		})
	if err != nil {
		println(err.Error())
	}
}
