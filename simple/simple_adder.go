package simple

import (
	"context"
	"os"

	chunker "github.com/ipfs/go-ipfs-chunker"
	"github.com/zzzhr1990/go-ipfs-util/datastores"

	// dag "github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs/importer/balanced"
	ihelper "github.com/ipfs/go-unixfs/importer/helpers"

	// "github.com/ipld/go-ipld-prime"
	bs "github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
)

const chnkr = "size-1048576"

const maxLinks = 1024
const rawLeaves = true

var cidBuilder = merkledag.V1CidPrefix()

// var liveCacheSize = uint64(256 << 10)

type SimpleAddResult struct {
	RootID      string
	WcsEtag     string
	SHA1        []byte
	HeadSHA1    []byte
	SpecialSHA1 []byte
	FileSize    int64
	NodeMap     map[string][]byte
}

func Add(ctx context.Context, path string) (*SimpleAddResult, error) {

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	//datast := datastores.NewHashDatastore(fileInfo.Size())
	//dstore := dssync.MutexWrap(datast)         // dssync.MutexWrap(ds.NewMapDatastore())
	//bstore := blockstore.NewBlockstore(dstore) // bserv.New(blockstore.NewBlockstore(ds), nil)
	bstore := datastores.NewEmptyBlockstore(fileInfo.Size(), ctx)
	bserv := bs.New(bstore, offline.Exchange(bstore))

	dag := merkledag.NewDAGService(bserv)
	chnk, err := chunker.FromString(reader, chnkr)
	if err != nil {
		return nil, err
	}
	params := ihelper.DagBuilderParams{
		Dagserv:    dag,
		RawLeaves:  rawLeaves,
		Maxlinks:   maxLinks,
		CidBuilder: cidBuilder,
		// Maxlinks: ihelper.DefaultLinksPerBlock,
	}

	db, err := params.New(chnk)
	if err != nil {
		return nil, err
	}

	nd, err := balanced.Layout(db)
	if err != nil {
		return nil, err
	}

	nodemap := bstore.NodeMap()
	result := &SimpleAddResult{
		RootID:      nd.Cid().String(),
		WcsEtag:     bstore.SumWcsEtag(),
		SHA1:        bstore.FileSHA1(),
		HeadSHA1:    bstore.HeadSHA1(),
		FileSize:    fileInfo.Size(),
		SpecialSHA1: bstore.SumXL(),
		NodeMap:     make(map[string][]byte, len(nodemap)),
	}

	for k, v := range nodemap {
		result.NodeMap[k] = v.RawData()
	}

	return result, nil
}
