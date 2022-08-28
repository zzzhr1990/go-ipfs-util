package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zzzhr1990/go-file-hasher/bthash"
	"github.com/zzzhr1990/go-ipfs-util/adder"
	"github.com/zzzhr1990/go-ipfs-util/datastores"
	"github.com/zzzhr1990/go-ipfs-util/simple"
)

func Xmain() {
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
	// fileNode.Links()
	for _, link := range fileNode.Links() {
		println(">")
		fmt.Println(link.Name)
		fmt.Println(link.Cid.String())
		// fmt.Println(link.c)
		// linkNode, err := fileNode.(context.Background(), link.Cid)
	}
	println("----------------------------------------------------")

	fmt.Println(fileNode.Cid().String())
	//cs, _ := hashcalc.CalcFileHash(filePath, 0)
	nd, r, err := simple.Add(filePath)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(nd.Cid().String())

	var ct int64 = 0
	var cs int = 0
	println("----------------------------------------------------")
	r.Reslove(nd.Cid().String(), &ct, &cs)

	println(ct, " fs: ", r.GetFileSize(), " block size: ", r.GetBlockCount())
}

func main() {
	if false {
		path := "/Users/herui/Downloads/image-PPC_M460EX-3.6.8012.img" // "/Users/herui/Downloads/2020年度武汉鲨鱼、鲨鱼北分总分合并报表及报告.pdf"
		_, is, err := simple.Add(path)
		if err != nil {
			println(err.Error())
		}
		xlHash1 := hex.EncodeToString(is.SumXL())
		xlHash2, err := datastores.CalcFileGcid(path)
		if err != nil {
			println(err.Error())
		}
		if xlHash1 != xlHash2 {
			println("not equal")
			println(xlHash1)
			println(xlHash2)
			println(path)
			// os.Exit(3)
		}
		etag, _ := datastores.GetEtag(path)
		if is.SumWcsEtag() != etag {
			println("not equal etag")
			println(etag)
			println(is.SumWcsEtag())
			println(path)
			os.Exit(3)
		}
		return
	}
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
				return nil
			}
			//fmt.Println(path, info.Size())
			if !info.IsDir() {
				nd, is, err := simple.Add(path)
				if err != nil {
					println(err.Error())
					return err
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
				hasher, _ := bthash.CreateNewHasher(path, -1, context.Background())
				if hasher.Sha1String() != hex.EncodeToString(is.FileSHA1()) {
					println("not equal sha1")
					println(hasher.Sha1String())
					println(hex.EncodeToString(is.FileSHA1()))
					println(path)
					os.Exit(3)
				}
				if hasher.HeadSha1String() != hex.EncodeToString(is.HeadSHA1()) {
					println("not equal head sha1")
					println(hasher.Sha1String())
					println(hex.EncodeToString(is.HeadSHA1()))
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
