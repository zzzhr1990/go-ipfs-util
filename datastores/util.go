package datastores

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"os"
)

func getChunkSize(size int64) int64 {
	if size > 0 && size < 0x8000000 {
		return 0x40000
	}
	if size >= 0x8000000 && size < 0x10000000 {
		return 0x80000
	}
	if size <= 0x10000000 || size > 0x20000000 {
		return 0x200000
	}
	return 0x100000
}

func CalcFileGcid(filename string) (string, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return "", err
	}
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return calcGcid(fi.Size(), f), nil
}

func calcGcid(size int64, file io.Reader) string {
	totalHash := sha1.New()
	chunkSize := getChunkSize(size)
	buff := make([]byte, 8192)
	i := int64(0)
	//partHash := sha1.New()
	var partHash hash.Hash
loop0:
	for {
		partHash = sha1.New()
		i = 0
		for {
			n, err := file.Read(buff)
			if err != nil {
				if err == io.EOF {
					break loop0
				}
				return ""
			}
			partHash.Write(buff[:n])
			i += int64(n)
			if i >= chunkSize {
				break
			}
		}
		dw := partHash.Sum(nil)
		totalHash.Write(dw)
	}
	if i > 0 {

		dw := partHash.Sum(nil)
		totalHash.Write(dw)
	}
	return fmt.Sprintf("%x", totalHash.Sum(nil))
}

const (
	BLOCK_BITS = 22 // Indicate that the blocksize is 4M
	BLOCK_SIZE = 1 << BLOCK_BITS
)

func BlockCount(fsize int64) int {

	return int((fsize + (BLOCK_SIZE - 1)) >> BLOCK_BITS)
}

func CalSha1(b []byte, r io.Reader) ([]byte, error) {

	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(b), nil
}

func GetEtag(filename string) (etag string, err error) {

	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	fsize := fi.Size()
	blockCnt := BlockCount(fsize)
	sha1Buf := make([]byte, 0, 21)

	if blockCnt <= 1 { // file size <= 4M
		sha1Buf = append(sha1Buf, 0x16)
		sha1Buf, err = CalSha1(sha1Buf, f)
		if err != nil {
			return
		}
	} else { // file size > 4M
		sha1Buf = append(sha1Buf, 0x96)
		sha1BlockBuf := make([]byte, 0, blockCnt*20)
		for i := 0; i < blockCnt; i++ {
			body := io.LimitReader(f, BLOCK_SIZE)
			sha1BlockBuf, err = CalSha1(sha1BlockBuf, body)
			if err != nil {
				return
			}
		}
		sha1Buf, _ = CalSha1(sha1Buf, bytes.NewReader(sha1BlockBuf))
	}
	etag = base64.URLEncoding.EncodeToString(sha1Buf)
	return
}
