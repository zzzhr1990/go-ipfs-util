package main

import "github.com/zzzhr1990/go-ipfs-util/calc"

func main() {
	filePath := "Top.Gunner.Danger.Zone.2022.1080p.WEBRip.x265-RARBG.mp4"
	file, err := calc.CalcFileHash(filePath)
	if err != nil {
		println(err)
		return
	} else {
		println(file)
	}
}
