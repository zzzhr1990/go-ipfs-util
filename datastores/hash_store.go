package datastores

/*
import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"os"
	"sync"

	"github.com/ipfs/go-datastore"
	dsq "github.com/ipfs/go-datastore/query"
)

// HashDatastore stores nothing, but conforms to the API.
// Useful to test with.
type HashDatastore struct {
	md5Hash        hash.Hash
	sha1Hash       hash.Hash
	headerSha1Hash hash.Hash
	totalXLHash    hash.Hash
	partXLHash     hash.Hash
	headerLeft     int
	chunkSize      int64
	chunkRead      int64
	// partReadCount  int
	//writeCount     int
	cachedXLHash []byte
	mutex        sync.Mutex
	fileSize     int64
	tl           int64
}

// NewHashDatastore constructs a null datastoe
func NewHashDatastore(fileSize int64) *HashDatastore {
	// sha256.New()
	return &HashDatastore{
		md5Hash:        md5.New(),
		sha1Hash:       sha1.New(),
		headerSha1Hash: sha1.New(),
		totalXLHash:    sha1.New(),
		partXLHash:     sha1.New(),
		chunkSize:      getChunkSize(fileSize),
		chunkRead:      0,
		fileSize:       fileSize,
	}
}

// Put implements Datastore.Put
func (d *HashDatastore) Put(ctx context.Context, key datastore.Key, value []byte) (err error) {
	d.tl += int64(len(value))
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		d.md5Hash.Write(value)
		wg.Done()
	}()
	go func() {
		d.sha1Hash.Write(value)
		wg.Done()
	}()
	wg.Wait()

	if d.headerLeft > 0 {
		//return
		bufferSize := len(value)
		needRead := d.headerLeft
		if needRead > bufferSize {
			needRead = bufferSize
		}
		d.headerSha1Hash.Write(value[:needRead])
		d.headerLeft -= needRead
	}

	dataToRead := value[:]
	for {

		if d.chunkRead >= d.chunkSize {
			d.chunkRead = 0
		}
		if d.chunkRead == 0 {
			d.partXLHash.Reset()
		}
		bufferSize := int64(len(dataToRead))

		needRead := d.chunkSize - d.chunkRead
		if needRead > bufferSize {
			needRead = bufferSize
		}
		d.partXLHash.Write(dataToRead[:needRead])
		dataToRead = dataToRead[needRead:]
		d.chunkRead += needRead
		if d.chunkRead == d.chunkSize {
			dw := d.partXLHash.Sum(nil)
			d.totalXLHash.Write(dw)
			// d.writeCount++
		}
		if len(dataToRead) == 0 {
			break
		}

		//if d.chunkRead == bufferSize {
		//	break
		//}
	}

	return nil
}

// Sync implements Datastore.Sync
func (d *HashDatastore) Sync(ctx context.Context, prefix datastore.Key) error {
	return nil
}

// Get implements Datastore.Get
func (d *HashDatastore) Get(ctx context.Context, key datastore.Key) (value []byte, err error) {
	return nil, datastore.ErrNotFound
}

// Has implements Datastore.Has
func (d *HashDatastore) Has(ctx context.Context, key datastore.Key) (exists bool, err error) {
	return false, nil
}

// Has implements Datastore.GetSize
func (d *HashDatastore) GetSize(ctx context.Context, key datastore.Key) (size int, err error) {
	return -1, datastore.ErrNotFound
}

// Delete implements Datastore.Delete
func (d *HashDatastore) Delete(ctx context.Context, key datastore.Key) (err error) {
	return nil
}

func (d *HashDatastore) Scrub(ctx context.Context) error {
	return nil
}

func (d *HashDatastore) Check(ctx context.Context) error {
	return nil
}

// Query implements Datastore.Query
func (d *HashDatastore) Query(ctx context.Context, q dsq.Query) (dsq.Results, error) {
	return dsq.ResultsWithEntries(q, nil), nil
}

func (d *HashDatastore) Batch(ctx context.Context) (datastore.Batch, error) {
	return datastore.NewBasicBatch(d), nil
}

func (d *HashDatastore) CollectGarbage(ctx context.Context) error {
	return nil
}

func (d *HashDatastore) DiskUsage(ctx context.Context) (uint64, error) {
	return 0, nil
}

func (d *HashDatastore) Close() error {
	return nil
}

func (d *HashDatastore) NewTransaction(ctx context.Context, readOnly bool) (datastore.Txn, error) {
	return &nullTxn{}, nil
}

type nullTxn struct{}

func (t *nullTxn) Get(ctx context.Context, key datastore.Key) (value []byte, err error) {
	return nil, nil
}

func (t *nullTxn) Has(ctx context.Context, key datastore.Key) (exists bool, err error) {
	return false, nil
}

func (t *nullTxn) GetSize(ctx context.Context, key datastore.Key) (size int, err error) {
	return 0, nil
}

func (t *nullTxn) Query(ctx context.Context, q dsq.Query) (dsq.Results, error) {
	return dsq.ResultsWithEntries(q, nil), nil
}

func (t *nullTxn) Put(ctx context.Context, key datastore.Key, value []byte) error {
	return nil
}

func (t *nullTxn) Delete(ctx context.Context, key datastore.Key) error {
	return nil
}

func (t *nullTxn) Commit(ctx context.Context) error {
	return nil
}

func (t *nullTxn) Discard(ctx context.Context) {}

*/
