package datastores

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"hash"
	"sync"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	legacy "github.com/ipfs/go-ipld-legacy"
	mdag "github.com/ipfs/go-merkledag"
	unixfs "github.com/ipfs/go-unixfs"
)

type EmptyBlockstore struct {
	wcsBlockHash     hash.Hash
	wcsChunkSize     int64
	wcsChunkRead     int64
	wcsBlockCount    int64
	wcsBlockHashList []byte
	sha1Hash         hash.Hash
	headerSha1Hash   hash.Hash
	totalXLHash      hash.Hash
	partXLHash       hash.Hash
	headerLeft       int
	chunkSize        int64
	chunkRead        int64
	// partReadCount  int
	//writeCount     int
	cachedXLHash       []byte
	cachedSHA1Hash     []byte
	cachedHeadSHA1Hash []byte
	cachedWcsEtag      string
	mutex              sync.Mutex
	fileSize           int64
	//countBlocksSize    int64
	nodeMap map[string]blocks.Block
	ctx     context.Context
}

func NewEmptyBlockstore(fileSize int64, ctx context.Context) *EmptyBlockstore {
	var wcsBlockSize int64 = 4 * 1024 * 1024
	wcsBlockCount := (fileSize + wcsBlockSize - 1) / wcsBlockSize // round up (5+2-1)/2=3, (4+2-1)/2=2 (0+2-1)/2=0 (1+2-1)/2=1
	return &EmptyBlockstore{
		wcsBlockHash:     sha1.New(),
		sha1Hash:         sha1.New(),
		headerSha1Hash:   sha1.New(),
		totalXLHash:      sha1.New(),
		partXLHash:       sha1.New(),
		chunkSize:        getChunkSize(fileSize),
		chunkRead:        0,
		fileSize:         fileSize,
		headerLeft:       128 * 1024,
		wcsChunkSize:     wcsBlockSize,
		wcsBlockCount:    wcsBlockCount,
		wcsBlockHashList: make([]byte, 0, wcsBlockCount*sha1.Size),
		nodeMap:          make(map[string]blocks.Block),
		// blockIDs:         make([]string, 0, (fileSize+1048576-1)/1048576),
		ctx: ctx,
	}
}

func (s *EmptyBlockstore) DeleteBlock(context.Context, cid.Cid) error {
	return nil
}
func (s *EmptyBlockstore) Has(context.Context, cid.Cid) (bool, error) {
	return false, nil
}
func (s *EmptyBlockstore) Get(context.Context, cid.Cid) (blocks.Block, error) {
	return nil, nil
}

// GetSize returns the CIDs mapped BlockSize
func (s *EmptyBlockstore) GetSize(context.Context, cid.Cid) (int, error) {
	return 0, nil
}

// Put puts a given block to the underlying datastore
func (s *EmptyBlockstore) Put(ctx context.Context, b blocks.Block) error {
	// context.WithCancel(ctx)
	cctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	value := b.RawData()
	if b.Cid().Prefix().Codec == cid.Raw {
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			dataToRead := value[:]
			for {
				if s.wcsChunkRead >= s.wcsChunkSize {
					s.wcsChunkRead = 0
				}
				if s.wcsChunkRead == 0 && s.wcsBlockCount > 1 {
					s.wcsBlockHash.Reset()
				}
				bufferSize := int64(len(dataToRead))

				needRead := s.wcsChunkSize - s.wcsChunkRead
				if needRead > bufferSize {
					needRead = bufferSize
				}
				s.wcsBlockHash.Write(dataToRead[:needRead])
				dataToRead = dataToRead[needRead:]
				s.wcsChunkRead += needRead
				if needRead == 0 {
					break
				}
				if s.wcsChunkRead == s.wcsChunkSize {
					//dw := s.partXLHash.Sum(nil)
					//s.totalXLHash.Write(dw)
					// s.writeCount++
					if s.wcsBlockCount > 1 {
						s.wcsBlockHashList = s.wcsBlockHash.Sum(s.wcsBlockHashList)
					}
				}

				//if s.chunkRead == bufferSize {
				//	break
				//}
			}
			wg.Done()
		}()
		go func() {
			s.sha1Hash.Write(value)
			wg.Done()
		}()
		go func() {
			dataToRead := value[:]
			for {

				if s.chunkRead >= s.chunkSize {
					s.chunkRead = 0
				}
				if s.chunkRead == 0 {
					s.partXLHash.Reset()
				}
				bufferSize := int64(len(dataToRead))

				needRead := s.chunkSize - s.chunkRead
				if needRead > bufferSize {
					needRead = bufferSize
				}
				s.partXLHash.Write(dataToRead[:needRead])
				dataToRead = dataToRead[needRead:]
				s.chunkRead += needRead
				if needRead == 0 {
					break
				}
				if s.chunkRead == s.chunkSize {
					dw := s.partXLHash.Sum(nil)
					s.totalXLHash.Write(dw)
					// s.writeCount++
				}

				//if s.chunkRead == bufferSize {
				//	break
				//}
			}
			wg.Done()
		}()
		ch := make(chan struct{}, 1)
		go func() {
			wg.Wait()
			ch <- struct{}{}
		}()

		select {
		case <-ch:
			break
		case <-cctx.Done():
			wg.Wait()
			return s.ctx.Err()
		}

		if s.headerLeft > 0 {
			//return
			bufferSize := len(value)
			needRead := s.headerLeft
			if needRead > bufferSize {
				needRead = bufferSize
			}
			s.headerSha1Hash.Write(value[:needRead])
			s.headerLeft -= needRead
		}
		//s.countBlocksSize += int64(len(b.RawData()))
		// s.blockIDs = append(s.blockIDs, b.Cid().String())
	} else {

		s.nodeMap[b.Cid().String()] = b
	}
	return nil
}

func (s *EmptyBlockstore) Reslove(cidString string, intptr *int64, currentBlockIDPtr *int) error {
	data, ok := s.nodeMap[cidString]
	if ok {
		_, decoded, err := DecodeBlock(context.Background(), data)
		if err != nil {
			return err
		}
		for _, link := range decoded.Links() {
			_, ok := s.nodeMap[link.Cid.String()]
			if !ok {
				*intptr += int64(link.Size)
				//currentID := *currentBlockIDPtr
				//if s.blockIDs[currentID] != link.Cid.String() {
				//	panic("not equal")
				//}
				*currentBlockIDPtr++
			}
			err := s.Reslove(link.Cid.String(), intptr, currentBlockIDPtr)
			if err != nil {
				return err
			}

			// link.Size
		}
	}
	return nil
}

func (s *EmptyBlockstore) GetBlockCount() int64 {
	i := int64(0)
	for _, v := range s.nodeMap {
		i += int64(len(v.RawData()))
	}
	return i
}

func (s *EmptyBlockstore) GetFileSize() int64 {

	return s.fileSize
}

// PutMany puts a slice of blocks at the same time using batching
// capabilities of the underlying datastore whenever possible.
func (s *EmptyBlockstore) PutMany(ctx context.Context, blocks []blocks.Block) error {
	for _, block := range blocks {
		err := s.Put(ctx, block)
		if err != nil {
			return err
		}
	}
	return nil
}

// AllKeysChan returns a channel from which
// the CIDs in the Blockstore can be read. It should respect
// the given context, closing the channel if it becomes Done.
func (s *EmptyBlockstore) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	ch := make(chan cid.Cid)
	close(ch)
	return ch, nil
}

// HashOnRead specifies if every read block should be
// rehashed to make sure it matches its CID.
func (s *EmptyBlockstore) HashOnRead(enabled bool) {

}

func DecodeBlock(ctx context.Context, b blocks.Block) (int64, ipld.Node, error) {
	nd, err := legacy.DecodeNode(ctx, b)
	if err != nil {
		return 0, nd, err
	}
	var size int64
	switch n := nd.(type) {
	case *mdag.RawNode:
		size = int64(len(n.RawData()))

	case *mdag.ProtoNode:
		fsNode, err := unixfs.FSNodeFromBytes(n.Data())
		if err != nil {
			panic(err)
		}

		switch fsNode.Type() {
		case unixfs.TFile, unixfs.TRaw:
			size = int64(fsNode.FileSize())

		case unixfs.TDirectory, unixfs.THAMTShard:
			// Dont allow reading directories
			panic("Cannot write to a directory")

		case unixfs.TMetadata:
			panic("Cannot write to a metadata file")
		case unixfs.TSymlink:
			panic("Cannot write to a symlink")
		default:
			panic("Unknown node type")
		}
	default:
		panic("Unknown node type")
	}
	return size, nd, nil
}

func (s *EmptyBlockstore) SumXL() []byte {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.cachedXLHash) > 0 {
		// return cachedXLHash
		return s.cachedXLHash
	}
	if s.chunkRead > 0 {

		dw := s.partXLHash.Sum(nil)
		s.totalXLHash.Write(dw)
		//d.totalXLHash.Write(d.partXLHash.Sum(nil))
		s.chunkRead = 0
	}

	s.cachedXLHash = s.totalXLHash.Sum(nil)

	return s.cachedXLHash
}

func (s *EmptyBlockstore) NodeMap() map[string]blocks.Block {
	return s.nodeMap
}

func (s *EmptyBlockstore) SumWcsEtag() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.cachedWcsEtag) > 0 {
		// return cachedXLHash
		return s.cachedWcsEtag
	}
	sha1Buf := make([]byte, 0, sha1.Size+1)
	if s.wcsBlockCount <= 1 {
		sha1Buf = append(sha1Buf, 0x16)
		sha1Buf = s.wcsBlockHash.Sum(sha1Buf)
	} else {
		if s.wcsChunkRead > 0 {

			s.wcsBlockHashList = s.wcsBlockHash.Sum(s.wcsBlockHashList)
			s.wcsChunkRead = 0
		}
		sha1Buf = append(sha1Buf, 0x96)
		s.wcsBlockHash.Reset()
		s.wcsBlockHash.Write(s.wcsBlockHashList)
		sha1Buf = s.wcsBlockHash.Sum(sha1Buf)
	}

	s.cachedWcsEtag = base64.URLEncoding.EncodeToString(sha1Buf)

	return s.cachedWcsEtag
}

func (s *EmptyBlockstore) FileSHA1() []byte {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.cachedSHA1Hash) > 0 {
		// return cachedXLHash
		return s.cachedSHA1Hash
	}

	s.cachedSHA1Hash = s.sha1Hash.Sum(nil)

	return s.cachedSHA1Hash
}

func (s *EmptyBlockstore) HeadSHA1() []byte {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.cachedHeadSHA1Hash) > 0 {
		// return cachedXLHash
		return s.cachedHeadSHA1Hash
	}

	s.cachedHeadSHA1Hash = s.headerSha1Hash.Sum(nil)

	return s.cachedHeadSHA1Hash
}
