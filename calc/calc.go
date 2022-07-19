package calc

import (
	"context"
	"errors"
	"fmt"

	// "log"
	"os"
	"strings"

	bs "github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-mfs"

	// cmds "github.com/ipfs/go-ipfs-cmds"
	files "github.com/ipfs/go-ipfs-files"

	cidutil "github.com/ipfs/go-cidutil"
	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-ipfs-pinner/dspinner"
	"github.com/ipfs/go-ipfs/core/coreunix"
	mdag "github.com/ipfs/go-merkledag"
	dagtest "github.com/ipfs/go-merkledag/test"
	ft "github.com/ipfs/go-unixfs"
	"github.com/ipfs/interface-go-ipfs-core/options"
	mh "github.com/multiformats/go-multihash"
)

type AddEvent struct {
	Name  string
	Hash  string `json:",omitempty"`
	Bytes int64  `json:",omitempty"`
	Size  string `json:",omitempty"`
}

func CalcFileHash(file string) (string, error) {

	// file := "/Volumes/Code/Users/zzzhr/Downloads/心电图一本通 2015.pdf"
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hashFunStr := "sha2-256"
	chunker := "size-262144"
	dopin := false
	hash := true
	fscache := false
	progress := false
	silent := false

	hashFunCode := mh.Names[strings.ToLower(hashFunStr)]
	opts := []options.UnixfsAddOption{

		options.Unixfs.Hash(hashFunCode),

		options.Unixfs.Inline(false),
		options.Unixfs.InlineLimit(32),

		options.Unixfs.Chunker(chunker),

		options.Unixfs.Pin(dopin),
		options.Unixfs.HashOnly(hash),
		options.Unixfs.FsCache(fscache),
		options.Unixfs.Nocopy(false),

		options.Unixfs.Progress(progress),
		options.Unixfs.Silent(silent),
	}

	//opts = append(opts, nil) // events option placeholder

	//events := make(chan interface{}, 32)
	//opts[len(opts)-1] = options.Unixfs.Events(events)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dstore := dssync.MutexWrap(ds.NewMapDatastore())
	bstore := blockstore.NewBlockstore(dstore)
	bserv := bs.New(bstore, offline.Exchange(bstore))

	dserv := mdag.NewDAGService(bserv)

	p, _ := dspinner.New(ctx, dstore, dserv)

	fileAdder, err := coreunix.NewAdder(ctx, p, blockstore.NewGCLocker(), dserv)
	if err != nil {
		return "", err
	}

	settings, prefix, err := options.UnixfsAddOptions(opts...)
	if err != nil {
		return "", err
	}

	fileAdder.Chunker = settings.Chunker
	if settings.Events != nil {
		fileAdder.Out = settings.Events
		fileAdder.Progress = settings.Progress
	}
	fileAdder.Pin = settings.Pin && !settings.OnlyHash
	fileAdder.Silent = settings.Silent
	fileAdder.RawLeaves = settings.RawLeaves
	fileAdder.NoCopy = settings.NoCopy
	fileAdder.CidBuilder = prefix

	switch settings.Layout {
	case options.BalancedLayout:
		// Default
	case options.TrickleLayout:
		fileAdder.Trickle = true
	default:
		// log.Fatalf("unknown layout: %d", settings.Layout)
		return "", errors.New(fmt.Sprintf("unknown layout: %d", settings.Layout))
	}

	if settings.Inline {
		fileAdder.CidBuilder = cidutil.InlineBuilder{
			Builder: fileAdder.CidBuilder,
			Limit:   settings.InlineLimit,
		}
	}

	//if settings.OnlyHash {
	md := dagtest.Mock()
	emptyDirNode := ft.EmptyDirNode()
	// Use the same prefix for the "empty" MFS root as for the file adder.
	emptyDirNode.SetCidBuilder(fileAdder.CidBuilder)
	mr, err := mfs.NewRoot(ctx, md, emptyDirNode, nil)
	if err != nil {
		// log.Fatalf("can't create new root: %v", err)
		return "", err
	}

	fileAdder.SetMfsRoot(mr)
	//}

	nd, err := fileAdder.AddAllAndPin(ctx, files.NewReaderFile(f))
	if err != nil {
		// log.Fatalf("can't create new root -- NDDD: %v", err)
		return "", err
	}

	return nd.Cid().String(), nil
}

// syncDagService is used by the Adder to ensure blocks get persisted to the underlying datastore
