package adder

/*
import (
	"context"
	"os"

	bs "github.com/ipfs/go-blockservice"
	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	w3fs "github.com/zzzhr1990/go-ipfs-util/fs"
)

// const targetChunkSize = 1024 * 1024 * 10

func HashFile(ctx context.Context, filename string, path string) (ipld.Node, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	dstore := dssync.MutexWrap(ds.NewNullDatastore()) // dssync.MutexWrap(ds.NewMapDatastore())
	bstore := blockstore.NewBlockstore(dstore)        // bserv.New(blockstore.NewBlockstore(ds), nil)
	bserv := bs.New(bstore, offline.Exchange(bstore))

	dag := merkledag.NewDAGService(bserv)
	dagFmtr, err := NewAdder(ctx, dag)
	if err != nil {
		return nil, err
	}

	root, err := dagFmtr.Add(filename, file, path, &w3fs.OsFs{})
	if err != nil {
		return nil, err
	}

	// If file is a dir, do not wrap in another.
	if info.IsDir() {
		mr, err := dagFmtr.MfsRoot()
		if err != nil {
			return nil, err
		}
		rdir := mr.GetDirectory()
		cdir, err := rdir.Child(info.Name())
		if err != nil {
			return nil, err
		}
		cnode, err := cdir.GetNode()
		if err != nil {
			return nil, err
		}
		root = cnode
	}
	/*
		//for _, link := range root.Links() {
		carReader, carWriter := io.Pipe()
		hash1 := md5.New()
		go func() {
			err = car.WriteCar(ctx, dag, []cid.Cid{root.Cid()}, carWriter)
			if err != nil {
				carWriter.CloseWithError(err)
				return
			}
			carWriter.Close()
		}()
		putCar(ctx, carReader, hash1)
		//}
*/
/*
	return root, nil
}

// PutCar uploads a CAR (Content Addressable Archive) to Web3.Storage.
/*
func putCar(ctx context.Context, car io.Reader, cryptoHash hash.Hash) (cid.Cid, error) {
	spltr, err := carbites.Split(car, targetChunkSize, carbites.Treewalk)
	if err != nil {
		return cid.Undef, err
	}

	var root cid.Cid
	for {
		r, err := spltr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return cid.Undef, err
		}
		size, err := io.Copy(cryptoHash, r)
		println("size: ", size)
		// TODO: concurrency

		if err != nil {
			return cid.Undef, err
		}
		//root = c
	}

	return root, nil
}

/*

func putCar2(ctx context.Context, car io.Reader) (cid.Cid, error) {
	carChunks := make(chan io.Reader)

	var root cid.Cid
	var wg sync.WaitGroup
	wg.Add(1)

	var sendErr error
	go func() {
		defer wg.Done()
		for r := range carChunks {
			// TODO: concurrency
			c, err := c.sendCar(ctx, r)
			if err != nil {
				sendErr = err
				break
			}
			root = c
		}
	}()

	err := carbites.Split(ctx, car, targetChunkSize, carbites.Treewalk, carChunks)
	if err != nil {
		return cid.Undef, err
	}
	wg.Wait()

	return root, sendErr
}
*/
